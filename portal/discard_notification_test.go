package portal_test

import (
	"testing"

	"github.com/nvvers/orbit-portal-client-golang/orbit"
	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_DiscardNotification(t *testing.T) {
	pc, err := portaltestcontainer.New()
	if err != nil {
		t.Fatalf("Failed to create test container: %v", err)
	}

	err = pc.Start(t.Context())
	if err != nil {
		t.Fatalf("Failed to start test container: %v", err)
	}
	defer pc.Stop(t.Context())

	client, err := pc.GetClient(t.Context())
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	testPoolName := "test-pool"
	err = client.SavePool(testPoolName)
	if err != nil {
		t.Fatalf("Failed to save pool: %v", err)
	}

	t.Run("discard notification", func(t *testing.T) {
		err := client.NotifyEvent("test-source", "test-subject", "test-type", map[string]string{"key": "value"}, nil)
		if err != nil {
			t.Fatalf("Failed to notify event: %v", err)
		}

		n, err := client.PullNotification(testPoolName)
		if err != nil {
			t.Fatalf("Failed to pull notification: %v", err)
		}

		err = client.DiscardNotification(testPoolName, n.Event.ID, nil)
		if err != nil {
			t.Fatalf("Failed to discard notification: %v", err)
		}
	})

	t.Run("discard notification with subsequent events", func(t *testing.T) {
		err := client.NotifyEvent("test-source", "test-subject", "test-type", map[string]string{"key": "value"}, nil)
		if err != nil {
			t.Fatalf("Failed to notify event: %v", err)
		}

		n, err := client.PullNotification(testPoolName)
		if err != nil {
			t.Fatalf("Failed to pull notification: %v", err)
		}

		err = client.DiscardNotification(testPoolName, n.Event.ID, []orbit.Event{{
			Source:  "test-source",
			Subject: "test-subject-2",
			Type:    "test-type-2",
			Data:    map[string]string{"key": "value"},
		}})
		if err != nil {
			t.Fatalf("Failed to discard notification with subsequent events: %v", err)
		}

		n2, err := client.PullNotification(testPoolName)
		if err != nil {
			t.Fatalf("Failed to pull subsequent notification: %v", err)
		}

		if n2.Event.Subject != "test-subject-2" {
			t.Fatalf("Expected subsequent event with subject 'test-subject-2', got '%s'", n2.Event.Subject)
		}
	})
}
