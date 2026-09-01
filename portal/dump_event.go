package portal

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nvvers/orbit-portal-client-golang/orbit"
)

func (c *Client) DumpEvent(dstFile string, er orbit.EventRecord) error {
	fs, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer fs.Close()

	zipWriter := zip.NewWriter(fs)
	defer zipWriter.Close()

	// ---

	for i, att := range er.Attachments {
		attWriter, err := zipWriter.Create("attachments/" + att.ID.String())
		if err != nil {
			return fmt.Errorf("failed to create attachment writer for %d: %v", i, err)
		}

		err = c.GetAttachment(att.ID, attWriter)
		if err != nil {
			return fmt.Errorf("failed to get attachment %d: %v", i, err)
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

	return nil
}
