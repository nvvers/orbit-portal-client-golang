package portal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_SaveAttachment(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "test-att.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

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

	f, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer f.Close()

	att, err := client.SaveAttachment(f)
	if err != nil {
		t.Fatalf("Failed to save attachment: %v", err)
	}

	if att.ID == uuid.Nil {
		t.Fatalf("Attachment ID is empty")
	}

	if att.Size == 0 {
		t.Fatalf("Attachment size is zero")
	}
}
