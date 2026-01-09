package portal

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/nvvers/orbit-portal-client-golang/internal/orbit"
)

func (c *Client) DumpEvent(dstFile string, er orbit.EventRecord) error {
	zipBuf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuf)

	// ---

	for i, att := range er.Attachments {
		attBuffer := new(bytes.Buffer)
		err := c.GetAttachment(att.ID, attBuffer)
		if err != nil {
			return fmt.Errorf("failed to get attachment %d: %v", i, err)
		}

		attWriter, err := zipWriter.Create("attachments/" + att.ID.String())
		if err != nil {
			return fmt.Errorf("failed to create attachment writer for %d: %v", i, err)
		}

		_, err = io.Copy(attWriter, attBuffer)
		if err != nil {
			return fmt.Errorf("failed to write attachment %d: %v", i, err)
		}
	}

	// ---

	eventJson, err := json.Marshal(er)
	if err != nil {
		return fmt.Errorf("failed to marshal event record: %v", err)
	}

	eventJsonWriter, err := zipWriter.Create("event.json")
	if err != nil {
		return fmt.Errorf("failed to create event.json writer: %v", err)
	}

	_, err = eventJsonWriter.Write(eventJson)
	if err != nil {
		return fmt.Errorf("failed to write event.json: %v", err)
	}

	// ---

	err = zipWriter.Close()
	if err != nil {
		return fmt.Errorf("failed to close zip writer: %v", err)
	}

	fs, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer fs.Close()

	_, err = io.Copy(fs, zipBuf)
	if err != nil {
		return fmt.Errorf("failed to write zip to file: %v", err)
	}

	return nil
}
