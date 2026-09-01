package portal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nvvers/orbit-portal-client-golang/orbit"
	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_GetAttachment(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "test-att.txt")
	testFileOut := filepath.Join(t.TempDir(), "test-att-out.txt")

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

	err = client.NotifyEvent("test-source", "test-subject", "test-type", map[string]string{"key": "value123456789"}, []orbit.Attachment{att})
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

	w, err := os.Create(testFileOut)
	defer w.Close()
	err = client.GetAttachment(n.Event.Attachments[0].ID, w)
	if err != nil {
		t.Fatalf("Failed to get attachment: %v", err)
	}

	data1, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	data2, err := os.ReadFile(testFileOut)
	if err != nil {
		t.Fatalf("Failed to read test output file: %v", err)
	}
	if string(data1) != string(data2) {
		t.Fatalf("Attachment data does not match original file")
	}
}
