# Orbit Portal Client for Go

A Go client library for the **Orbit Portal** — an event distribution service that ingests events,
fans them out into named *pools* based on subject and type filters, and delivers them to consumers
as acknowledgeable *notifications*.

The library wraps the Portal HTTP API in a small client package. It covers publishing events, managing pools and JSON schemas,
consuming notifications (pull-based or streaming), handling binary attachments, and reading/writing
offline event dumps. It ships with a [testcontainers](https://golang.testcontainers.org/) helper so
you can spin up a real Portal instance in your own test suite.

## Table of contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Concepts](#concepts)
- [Creating a client](#creating-a-client)
- [Managing pools](#managing-pools)
- [Publishing events](#publishing-events)
- [Attachments](#attachments)
- [Receiving notifications](#receiving-notifications)
- [Completing a notification](#completing-a-notification)
- [Event schemas](#event-schemas)
- [Event dumps](#event-dumps)
- [Testing against a real portal](#testing-against-a-real-portal)
- [Project layout](#project-layout)
- [Known limitations](#known-limitations)
- [Contributing](#contributing)

## Requirements

- **Go 1.26** or later
- A reachable Orbit Portal instance and an API token
- Docker — only for running the integration tests (see [Testing against a real portal](#testing-against-a-real-portal))

## Installation

```bash
go get github.com/nvvers/orbit-portal-client-golang
```

The only non-test dependency is `github.com/google/uuid`; testcontainers is pulled in solely by the
`portaltestcontainer` package.

## Quick start

```go
package main

import (
	"fmt"
	"log"

	"github.com/nvvers/orbit-portal-client-golang/portal"
)

func main() {
	client, err := portal.NewClient("https://portal.example.com", portal.WithToken("my-token"))
	if err != nil {
		log.Fatal(err)
	}

	// Make sure the portal is reachable and the token is valid.
	if err := client.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := client.VerifyToken(); err != nil {
		log.Fatal(err)
	}

	// Create a pool that collects all "order.created" events.
	if err := client.SavePool("orders", portal.WithEventTypes([]string{"order.created"})); err != nil {
		log.Fatal(err)
	}

	// Publish an event.
	err = client.NotifyEvent(
		"checkout-service",              // source
		"order/4711",                    // subject
		"order.created",                 // type
		map[string]any{"total": 129.90}, // data
		nil,                             // attachments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Consume one notification and acknowledge it.
	n, err := client.PullNotification("orders")
	if err != nil {
		log.Fatal(err)
	}
	if n != nil {
		fmt.Printf("received %s (%s)\n", n.Event.Subject, n.Event.Type)

		if err := client.AcknowledgeNotification("orders", n.Event.ID, nil); err != nil {
			log.Fatal(err)
		}
	}
}
```

## Concepts

| Concept | Description |
| --- | --- |
| **Event** | An immutable fact published to the Portal. Carries `source`, `subject`, `type`, an arbitrary JSON `data` payload and optional attachments. |
| **Pool** | A named, server-side subscription. A pool declares which events it is interested in (by event type, subject pattern and filter expression) and turns matching events into notifications. |
| **Notification** | A pool's delivery of one event to one consumer. It is redelivered until it is acknowledged, discarded or expires. |
| **Attachment** | A binary blob stored by the Portal and referenced by an event through its id. |
| **Schema** | A JSON schema registered per event type. The Portal validates incoming events against it. |
| **`.oed` file** | *Orbit Event Dump* — a zip archive containing `event.json` plus one entry per attachment under `attachments/<attachment-id>`. |

## Creating a client

```go
client, err := portal.NewClient(baseURL, opts...)
```

| Option | Description |
| --- | --- |
| `portal.WithToken(token)` | Sends `Authorization: Bearer <token>` with every request. Omit it for portals that allow anonymous access. |
| `portal.WithHttpClient(hc)` | Supplies your own `*http.Client` — use it to configure timeouts, proxies or TLS. Returns an error if `hc` is `nil`. |

All methods return wrapped errors; unexpected HTTP status codes are reported together with the
response body, so `errors.Is`/`%w` unwrapping works up the chain.

## Managing pools

```go
err := client.SavePool("orders",
	portal.WithSubjectPattern(`^order/.*`),
	portal.WithEventTypes([]string{"order.created", "order.cancelled"}),
	portal.WithFilterExpression(`data.total > 100`),
	portal.WithNotificationRetryInterval(30*time.Second),
	portal.WithNotificationTTL(7*24*time.Hour),
)

err = client.DeletePool("orders")
```

`SavePool` is an upsert — calling it again with the same name updates the existing pool.

| Pool option | Default | Description |
| --- | --- | --- |
| `WithSubjectPattern(pattern)` | `.*` | Regular expression matched against the event subject. |
| `WithEventTypes(types)` | `[]` (all types) | Restricts the pool to the given event types. |
| `WithFilterExpression(expr)` | *(none)* | Server-side expression evaluated against the event, applied on top of subject and type matching. |
| `WithNotificationRetryInterval(d)` | 5 minutes | How long an unacknowledged notification stays invisible before it is redelivered. Must be > 0. |
| `WithNotificationTTL(d)` | 24 hours | How long a notification survives before it expires. Must be > 0. |

## Publishing events

```go
err := client.NotifyEvent(source, subject, eventType, data, attachments)
```

`data` is marshalled to JSON, so any type accepted by `encoding/json` works. If a schema is
registered for `eventType`, the portal rejects events that do not validate against it.

## Attachments

```go
att, err := client.SaveAttachmentFromFile("invoice.pdf")
// or: att, err := client.SaveAttachment(reader)

err = client.NotifyEvent("billing", "invoice/2026-08", "invoice.issued", data, []orbit.Attachment{att})

var buf bytes.Buffer
err = client.GetAttachment(att.ID, &buf)
```

Attachments are uploaded first and referenced by ID from the event, so large payloads never travel
inside the event body. `GetAttachment` streams into any `io.Writer`.

## Receiving notifications

The client offers two delivery styles. Both hand out notifications one at a time and expect an
explicit disposition — a notification that is neither acknowledged, discarded nor resubmitted is
redelivered after the pool's retry interval.

**Pull** — fetches the next pending notification, or returns `(nil, nil)` when the pool is
empty. It is the right fit for cron-style workers and batch jobs.

```go
for {
	n, err := client.PullNotification("orders")
	if err != nil {
		return err
	}
	if n == nil {
		break // nothing pending
	}

	if err := handle(n); err != nil {
		// Try again in five minutes instead of failing the notification.
		if err := client.RequestResubmission("orders", n.Event.ID, time.Now().Add(5*time.Minute)); err != nil {
			return err
		}
		continue
	}

	if err := client.AcknowledgeNotification("orders", n.Event.ID, nil); err != nil {
		return err
	}
}
```

**Observe** — keeps a streaming connection open and invokes the callback for every
notification. Server heartbeats are filtered out transparently; the call returns `nil` when the
context is cancelled, when the server closes the stream, or when the connection ends.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := client.ObserveNotifications(ctx, "orders", func(n orbit.Notification) {
	if err := handle(n); err != nil {
		_ = client.DiscardNotification("orders", n.Event.ID, nil)
		return
	}
	_ = client.AcknowledgeNotification("orders", n.Event.ID, nil)
})
```

The callback runs synchronously on the receiving goroutine — the connection is only released for
the next notification after it returns, which keeps delivery strictly sequential. Do the heavy
lifting elsewhere if you need concurrency.

## Completing a notification

Every notification must be completed in one of three ways, otherwise it is redelivered:

| Method | Effect |
| --- | --- |
| `AcknowledgeNotification(pool, eventID, subsequentEvents)` | Marks the notification as processed. Any `subsequentEvents` are published atomically with the acknowledgement — the idiomatic way to chain a workflow without losing events on a crash. |
| `DiscardNotification(pool, eventID, subsequentEvents)` | Drops the notification without treating it as successfully processed, optionally publishing follow-up events. |
| `RequestResubmission(pool, eventID, resubmissionTime)` | Defers the notification until the given time. |

Passing `subsequentEvents` to acknowledge or discard lets you publish new events as part of
completing the notification, so processing and follow-up publication are a single operation.

## Event schemas

```go
err := client.SaveSchema("order.created", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"total": map[string]any{"type": "number"},
	},
	"required": []string{"total"},
})

err = client.DeleteSchema("order.created")
```

## Event dumps

Dumps are self-contained `.oed` ZIP archives holding `event.json` and every referenced attachment.
They are useful for reproducing production events in a test environment, for archiving, and for
offline replay.

```go
// Write a received event, including its attachments, to disk.
err := client.DumpEvent("event-4711.oed", n.Event)

// Build a dump from scratch, without a portal round-trip for the attachments.
err = client.DumpEventNotification("new-event.oed", "billing", "invoice/1", "invoice.issued", data,
	[]portal.AttachmentProvider{
		func() (orbit.Attachment, io.ReadCloser) {
			f, _ := os.Open("invoice.pdf")
			return orbit.Attachment{Name: "invoice.pdf"}, f
		},
	},
)

// Replay a dump into a portal — attachments are re-uploaded and their IDs rewritten.
err = client.NotifyEventFromDumpFile("event-4711.oed")

// Grep a folder of dumps; matches are printed to stdout.
err = client.Search("./dumps", "invoice/2026-08")
```

`AttachmentProvider` returns the attachment metadata plus a reader for its content; the reader is
closed by the client. A zero `ID` is filled in with a fresh UUID.

## Testing against a real portal

`portaltestcontainer` starts a disposable Portal instance in Docker, generates a matching
configuration and hands you a preconfigured client.

```go
func TestOrders(t *testing.T) {
	pc, err := portaltestcontainer.New(portaltestcontainer.WithToken("secret"))
	if err != nil {
		t.Fatal(err)
	}

	if err := pc.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	defer pc.Stop(t.Context())

	client, err := pc.GetClient(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if err := client.SavePool("orders"); err != nil {
		t.Fatal(err)
	}
	// ...
}
```

| Container option | Default | Description |
| --- | --- | --- |
| `WithImageTag(tag)` | `latest` | Tag of the `cr.nv-online.dev/nv/orbit/portal` image. |
| `WithToken(token)` | `secret` | Token of the generated test client, which is granted all rights on all pools. |
| `WithPort(port)` | `8080` | Container-internal API port. The host port is mapped dynamically — use `GetBaseURL`. |

Alongside `GetClient`, the container exposes `GetHost`, `GetMappedPort`, `GetBaseURL`, `Token`,
`IsRunning` and `Stop`.

Running the test suite requires a working Docker daemon and pull access to the Portal image:

```bash
go test ./...
```

## Project layout

```
portal/               public client — one file per API operation
portaltestcontainer/  testcontainers helper for integration tests
internal/orbit/       domain types (Event, EventRecord, Notification, Attachment, Message)
internal/portalapi/   request/response DTOs for the HTTP API
```

## Known limitations

- **The domain types live under `internal/`.** `orbit.Event`, `orbit.Notification`,
  `orbit.EventRecord` and `orbit.Attachment` appear in the public signatures of `portal`, but Go's
  internal-package rule prevents other modules from importing them. Until these types move to a
  public package, external consumers cannot call `ObserveNotifications`, attach files to events, or
  pass subsequent events to acknowledge/discard. This must be resolved before the first tagged
  release.
- **Most methods do not take a `context.Context`.** Only `ObserveNotifications` is context-aware;
  the remaining operations use `context.Background()` internally, so they cannot be cancelled or
  deadline-bound by the caller. Pass a `*http.Client` with a `Timeout` via `WithHttpClient` as a
  stopgap.
- **`Search` prints to stdout** instead of returning results, which limits it to CLI-style use.

## Contributing

Issues and pull requests are welcome. Please make sure `go build ./...`, `go vet ./...` and
`go test ./...` pass before submitting; the test suite needs Docker.
