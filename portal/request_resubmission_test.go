package portal_test

import (
	"testing"
	"time"

	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_RequestResubmission(t *testing.T) {
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
		if n == nil {
			t.Fatalf("Expected a notification, but got nil")
		}

		n2, err := client.PullNotification(testPoolName)
		if err != nil {
			t.Fatalf("Failed to pull notification: %v", err)
		}
		if n2 != nil {
			t.Fatalf("Expected no notification, but got one: %v", n2)
		}

		err = client.RequestResubmission(testPoolName, n.Event.ID, time.Now().Add(1*time.Second))
		if err != nil {
			t.Fatalf("Failed to request resubmission: %v", err)
		}

		time.Sleep(2 * time.Second)

		n3, err := client.PullNotification(testPoolName)
		if err != nil {
			t.Fatalf("Failed to pull notification: %v", err)
		}
		if n3 == nil {
			t.Fatalf("Expected a notification, but got nil")
		}
	})
}
