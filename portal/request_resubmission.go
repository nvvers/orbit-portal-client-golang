package portal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
)

func (c *Client) RequestResubmission(poolName string, eventID uuid.UUID, resubmissionTime time.Time) error {
	req := portalapi.RequestResubmissionRequest{
		PoolName:         poolName,
		EventID:          eventID,
		ResubmissionTime: resubmissionTime,
	}

	hReq, err := c.createPostRequestWithJsonBody(context.Background(), "/api/v1/request-resubmission", nil, req)
	if err != nil {
		return fmt.Errorf("failed to request resubmission: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to request resubmission: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to request resubmission: %w", err)
	}

	return nil
}
