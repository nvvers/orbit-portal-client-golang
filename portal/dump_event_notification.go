package portal

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	orbit "github.com/nvvers/orbit-portal-client-golang/orbit"
)

type AttachmentProvider func() (orbit.Attachment, io.ReadCloser)

func (c *Client) DumpEventNotification(dstFile string, source, subject, eventType string, data any, attachments []AttachmentProvider) error {
	zipBuf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuf)

	e := orbit.Event{
		Source:  source,
		Subject: subject,
		Type:    eventType,
		Data:    data,
	}

	if e.Source == "" {
		e.Source = c.defaultSource
	}

	// --- Attachments ---

	for i, attProvider := range attachments {
		att, attReader := attProvider()
		err := func() error {
			defer attReader.Close()

			if att.ID == uuid.Nil {
				att.ID = uuid.New()
			}

			attWriter, err := zipWriter.Create("attachments/" + att.ID.String())
			if err != nil {
				return fmt.Errorf("failed to create attachment writer for %d: %v", i, err)
			}

			_, err = io.Copy(attWriter, attReader)
			if err != nil {
				return fmt.Errorf("failed to write attachment %d: %v", i, err)
			}

			e.Attachments = append(e.Attachments, att)
			return nil
		}()

		if err != nil {
			return err
		}
	}

	eventData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}

	// --- Event JSON ---

	er := orbit.EventRecord{
		ID:          uuid.New(),
		Time:        time.Now(),
		Source:      e.Source,
		Subject:     e.Subject,
		Type:        e.Type,
		Data:        eventData,
		Attachments: e.Attachments,
	}

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
