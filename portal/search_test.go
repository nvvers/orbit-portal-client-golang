package portal_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/nvvers/orbit-portal-client-golang/orbit"
	"github.com/nvvers/orbit-portal-client-golang/portal"
)

func TestClient_Search(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-att.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	client, err := portal.NewClient("") // Naked client for testing
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	f, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer f.Close()

	dumpFile := filepath.Join(tempDir, "test.oed")
	err = client.DumpEventNotification(dumpFile, "test-source", "test-subject", "test-type", map[string]string{"key": "value123456789"},
		[]portal.AttachmentProvider{
			func() (orbit.Attachment, io.ReadCloser) {
				return orbit.Attachment{
						ID:          uuid.UUID{},
						Name:        "test-att.txt",
						Size:        11,
						Description: "",
					},
					f
			},
		},
	)

	results, err := client.Search(tempDir, "value123456789")
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("No results found")
	}
}
