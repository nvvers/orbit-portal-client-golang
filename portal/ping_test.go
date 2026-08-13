package portal_test

import (
	"testing"

	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_Ping(t *testing.T) {
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

	t.Run("ping", func(t *testing.T) {
		err := client.Ping()
		if err != nil {
			t.Fatalf("Failed to notify event: %v", err)
		}
	})
}
