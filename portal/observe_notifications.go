package portal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nvvers/orbit-portal-client-golang/internal/orbit"
	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
)

func (c *Client) ObserveNotifications(ctx context.Context, poolName string, yield func(orbit.Notification)) error {
	req := portalapi.ObserveNotificationsRequest{PoolName: poolName}

	hReq, err := c.createPostRequestWithJsonBody(ctx, "/api/v1/observe-notifications", nil, req)
	if err != nil {
		return fmt.Errorf("failed to create request for observing notifications: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to observe notifications: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to observe notifications: %w", err)
	}

	connId := hRes.Header.Get("X-Orbit-Connection-ID")
	if connId == "" {
		return fmt.Errorf("missing X-Orbit-Connection-ID header in response")
	}

	decoder := json.NewDecoder(hRes.Body)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			var msg orbit.Message
			if err := decoder.Decode(&msg); errors.Is(err, io.EOF) {
				return nil
			} else if errors.Is(err, context.DeadlineExceeded) {
				return nil
			} else if err != nil {
				return fmt.Errorf("failed to unmarshal notification line: %w", err)
			}

			switch msg.Type {
			case orbit.MessageTypeHeartbeat:
				continue

			case orbit.MessageTypeNotification:
				var notification orbit.Notification
				if err := json.Unmarshal(msg.Payload, &notification); err != nil {
					return fmt.Errorf("failed to unmarshal notification payload: %w", err)
				}

				yield(notification)

				// Release connection, to allow the portal to send more notifications to this client.

				hReq, err := c.createPostRequestWithJsonBody(ctx, "/api/v1/observe-notifications", url.Values{"release": {connId}}, req)
				if err != nil {
					return fmt.Errorf("failed to release connection: %w", err)
				}

				hRes, err = c.httpClient.Do(hReq)
				if err != nil {
					return fmt.Errorf("failed to release connection: %w", err)
				}

				err = expect(hRes, http.StatusOK)
				if err != nil {
					hRes.Body.Close()
					return fmt.Errorf("failed to release connection: %w", err)
				}

				hRes.Body.Close()

			case orbit.MessageTypeClose:
				return nil
			}
		}
	}
}
