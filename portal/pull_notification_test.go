package portal_test

import (
	"testing"

	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_PullNotification(t *testing.T) {
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

	t.Run("pull notification", func(t *testing.T) {
		err := client.NotifyEvent("test-source", "test-subject", "test-type", map[string]string{"key": "value"}, nil)
		if err != nil {
			t.Fatalf("Failed to notify event: %v", err)
		}

		err = client.NotifyEvent("test-source", "test-subject2", "test-type", map[string]string{"key": "value"}, nil)
		if err != nil {
			t.Fatalf("Failed to notify event: %v", err)
		}

		n, err := client.PullNotification(testPoolName)
		if err != nil {
			t.Fatalf("Failed to pull notification: %v", err)
		}
		if n == nil {
			t.Fatalf("Failed to pull notification")
		}

		n, err = client.PullNotification(testPoolName)
		if err != nil {
			t.Fatalf("Failed to pull notification: %v", err)
		}
		if n == nil {
			t.Fatalf("Failed to pull notification")
		}

		n, err = client.PullNotification(testPoolName)
		if err != nil {
			t.Fatalf("Failed to pull notification: %v", err)
		}
		if n != nil {
			t.Fatalf("Expected no notification, but got one: %v", n)
		}
	})
}
