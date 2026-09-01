package portal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
	"github.com/nvvers/orbit-portal-client-golang/orbit"
)

func (c *Client) DiscardNotification(poolName string, eventID uuid.UUID, subsequentEvents []orbit.Event) error {
	req := portalapi.DiscardNotificationRequest{
		PoolName:         poolName,
		EventID:          eventID,
		SubsequentEvents: subsequentEvents,
	}

	hReq, err := c.createPostRequestWithJsonBody(context.Background(), "/api/v1/discard-notification", nil, req)
	if err != nil {
		return fmt.Errorf("failed to discard notification: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to discard notification: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to discard notification: %w", err)
	}

	return nil
}
