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

	"github.com/nvvers/orbit-portal-client-golang/orbit"
)

func (c *Client) Search(folder string, pattern string) ([]orbit.EventRecord, error) {
	var results []orbit.EventRecord

	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
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

					results = append(results, event)
				}

				_ = f.Close()
				break
			}

			return nil
		}()

		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
