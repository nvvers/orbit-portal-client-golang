package portal_test

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/nvvers/orbit-portal-client-golang/internal/orbit"
	"github.com/nvvers/orbit-portal-client-golang/portal"
)

func TestClient_DumpEventNotification(t *testing.T) {
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

	zipReader, err := zip.OpenReader(dumpFile)
	if err != nil {
		t.Fatalf("Failed to open dump file: %v", err)
	}
	defer zipReader.Close()

	var results []orbit.EventRecord
	for _, file := range zipReader.File {
		if file.Name != "event.json" {
			continue
		}

		func() {
			f, err := file.Open()
			if err != nil {
				t.Fatal("failed to open event file in dump", err)
			}
			defer f.Close()

			data, err := io.ReadAll(io.LimitReader(f, 1*1024*1024)) // limit to 1MB
			if err != nil {

				t.Fatal("failed to read event file in dump", err)
			}

			if strings.Contains(string(data), "value123456789") {
				var event orbit.EventRecord
				if err := json.NewDecoder(bytes.NewReader(data)).Decode(&event); err != nil {
					_ = f.Close()
					t.Fatal("failed to decode event from dump", err)
				}

				results = append(results, event)
			}

			_ = f.Close()
		}()
	}

	if len(results) == 0 {
		t.Fatalf("Failed to find event in dump file")
	}
}
