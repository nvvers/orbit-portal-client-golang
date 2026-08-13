package portal

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
)

func (c *Client) GetAttachment(attachmentID uuid.UUID, w io.Writer) error {
	req := portalapi.GetAttachmentRequest{AttachmentID: attachmentID}

	hReq, err := c.createPostRequestWithJsonBody(context.Background(), "/api/v1/get-attachment", nil, req)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to get attachment: %w", err)
	}

	if _, err := io.Copy(w, hRes.Body); err != nil {
		return fmt.Errorf("failed to write attachment to writer: %w", err)
	}

	return nil
}
