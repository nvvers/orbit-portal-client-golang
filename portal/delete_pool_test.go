package portal_test

import (
	"testing"

	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_DeletePool(t *testing.T) {
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

	err = client.DeletePool(testPoolName)
	if err != nil {
		t.Fatalf("Failed to delete pool: %v", err)
	}
}
