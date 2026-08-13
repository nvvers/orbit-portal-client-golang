package portal_test

import (
	"context"
	"testing"
	"time"

	"github.com/nvvers/orbit-portal-client-golang/internal/orbit"
	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_ObserveNotifications(t *testing.T) {
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

	t.Run("observe notifications", func(t *testing.T) {
		err := client.NotifyEvent("test-source", "test-subject", "test-type", map[string]string{"key": "value"}, nil)
		if err != nil {
			t.Fatalf("Failed to notify event: %v", err)
		}

		err = client.NotifyEvent("test-source", "test-subject2", "test-type", map[string]string{"key": "value"}, nil)
		if err != nil {
			t.Fatalf("Failed to notify event: %v", err)
		}

		nCount := 0
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()

		err = client.ObserveNotifications(ctx, testPoolName, func(n orbit.Notification) {
			nCount++
		})

		if err != nil {
			t.Fatalf("Failed to observe notifications: %v", err)
		}

		if nCount != 2 {
			t.Fatalf("Expected to observe 2 notifications, but observed %d", nCount)
		}
	})
}
