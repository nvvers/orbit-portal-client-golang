package portal_test

import (
	"testing"

	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_DeleteSchema(t *testing.T) {
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

	// Save JSON Schema
	err = client.SaveSchema("com.example.test-tested", map[string]any{})
	if err != nil {
		t.Fatalf("Failed to save schema: %v", err)
	}

	// Delete JSON Schema
	err = client.DeleteSchema("com.example.test-tested")
	if err != nil {
		t.Fatalf("Failed to delete schema: %v", err)
	}
}
