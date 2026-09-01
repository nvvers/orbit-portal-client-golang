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

	orbit2 "github.com/nvvers/orbit-portal-client-golang/orbit"
	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_DumpEvent(t *testing.T) {
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

	err = client.NotifyEvent("test-source", "test-subject", "test-type", map[string]string{"key": "value123456789"}, []orbit2.Attachment{att})
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

	dumpFile := filepath.Join(t.TempDir(), "test.oed")
	err = client.DumpEvent(dumpFile, n.Event)
	if err != nil {
		t.Fatalf("Failed to dump event: %v", err)
	}

	zipReader, err := zip.OpenReader(dumpFile)
	if err != nil {
		t.Fatalf("Failed to open dump file: %v", err)
	}
	defer zipReader.Close()

	var results []orbit2.EventRecord
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
				var event orbit2.EventRecord
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
