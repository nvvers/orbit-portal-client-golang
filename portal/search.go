package portal

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nvvers/orbit-portal-client-golang/internal/orbit"
)

func (c *Client) Search(folder string, pattern string) error {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.ToLower(filepath.Ext(entry.Name())) != ".oed" {
			continue
		}

		err := func() error {
			zipReader, err := zip.OpenReader(filepath.Join(folder, entry.Name()))
			if err != nil {
				return fmt.Errorf("failed to open dump file: %w", err)
			}
			defer zipReader.Close()

			for _, file := range zipReader.File {
				if file.Name != "event.json" {
					continue
				}

				f, err := file.Open()
				if err != nil {
					return fmt.Errorf("failed to open event file in dump: %w", err)
				}

				data, err := io.ReadAll(io.LimitReader(f, 1*1024*1024)) // limit to 1MB
				if err != nil {
					return fmt.Errorf("failed to read event file in dump: %w", err)
				}

				if strings.Contains(string(data), pattern) {
					var event orbit.EventRecord
					if err := json.NewDecoder(bytes.NewReader(data)).Decode(&event); err != nil {
						return fmt.Errorf("failed to decode event from dump: %w", err)
					}

					fmt.Printf(
						"File:      %s\n"+
							"EventID:   %s\n"+
							"Time:      %s\n"+
							"Event-Data:\n"+
							"  Source:  %s\n"+
							"  Subject: %s\n"+
							"  Type:    %s\n"+
							"  Data:    %s\n"+
							"  Attachments: %+v\n\n",
						entry.Name(), event.ID, event.Time,
						event.Source, event.Subject, event.Type, event.Data, event.Attachments,
					)
				}

				_ = f.Close()
				break
			}

			return nil
		}()

		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", entry.Name(), err)
		}
	}

	return nil
}
