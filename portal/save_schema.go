package portal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
)

func (c *Client) SaveSchema(eventType string, schema map[string]any) error {
	if eventType == "" {
		return fmt.Errorf("event type must not be empty")
	}

	req := portalapi.SaveSchemaRequest{
		EventType: eventType,
		Schema:    schema,
	}

	hReq, err := c.createPostRequestWithJsonBody(context.Background(), "/api/v1/save-schema", nil, req)
	if err != nil {
		return fmt.Errorf("failed to save schema: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to save schema: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to save schema: %w", err)
	}

	return nil
}
