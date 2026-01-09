package portal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
)

func (c *Client) DeleteSchema(eventType string) error {
	req := portalapi.DeleteSchemaRequest{
		EventType: eventType,
	}

	hReq, err := c.createPostRequest(context.Background(), "/api/v1/delete-schema", "", req)
	if err != nil {
		return fmt.Errorf("failed to delete schema: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to delete schema: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to delete schema: %w", err)
	}

	return nil
}
