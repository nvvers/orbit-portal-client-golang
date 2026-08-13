package portal

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nvvers/orbit-portal-client-golang/internal/orbit"
)

func (c *Client) NotifyEvent(source, subject, eventType string, data any, attachments []orbit.Attachment) error {
	req := orbit.Event{
		Source:      source,
		Subject:     subject,
		Type:        eventType,
		Data:        data,
		Attachments: attachments,
	}

	hReq, err := c.createPostRequestWithJsonBody(context.Background(), "/api/v1/notify-event", nil, req)
	if err != nil {
		return fmt.Errorf("failed to notify event: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to notify event: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to notify event: %w", err)
	}

	return nil
}

func (c *Client) NotifyEventFromDumpFile(dumpFile string) error {
	zipReader, err := zip.OpenReader(dumpFile)
	if err != nil {
		return fmt.Errorf("failed to open dump file: %w", err)
	}
	defer zipReader.Close()

	var event orbit.Event
	for _, file := range zipReader.File {
		if file.Name != "event.json" {
			continue
		}

		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open event file in dump: %w", err)
		}

		if err := json.NewDecoder(f).Decode(&event); err != nil {
			return fmt.Errorf("failed to decode event from dump: %w", err)
		}

		_ = f.Close()
		break
	}

	for i, att := range event.Attachments {
		for _, file := range zipReader.File {
			if file.Name != "attachments/"+att.ID.String() {
				continue
			}

			f, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open attachment file in dump: %w", err)
			}

			newAtt, err := c.SaveAttachment(f)
			if err != nil {
				return fmt.Errorf("failed to save attachment from dump: %w", err)
			}

			event.Attachments[i].ID = newAtt.ID
		}
	}

	return c.NotifyEvent(event.Source, event.Subject, event.Type, event.Data, event.Attachments)
}
