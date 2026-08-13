package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/nvvers/orbit-portal-client-golang/internal/orbit"
)

func (c *Client) SaveAttachmentFromFile(filepath string) (orbit.Attachment, error) {
	nilAttachment := orbit.Attachment{}

	f, err := os.Open(filepath)
	if err != nil {
		return nilAttachment, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer f.Close()

	return c.SaveAttachment(f)
}

func (c *Client) SaveAttachment(reader io.Reader) (orbit.Attachment, error) {
	nilAttachment := orbit.Attachment{}

	hReq, err := c.createPostRequest(context.Background(), "/api/v1/save-attachment", nil, reader)
	if err != nil {
		return nilAttachment, fmt.Errorf("failed to create request: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return nilAttachment, fmt.Errorf("failed to send request: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return nilAttachment, fmt.Errorf("failed to save attachment: %w", err)
	}

	var response orbit.Attachment
	if err := json.NewDecoder(hRes.Body).Decode(&response); err != nil {
		return nilAttachment, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.ID == uuid.Nil {
		return nilAttachment, fmt.Errorf("received empty attachment ID")
	}

	return response, nil
}
