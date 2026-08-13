package portal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
)

func (c *Client) DeletePool(poolName string) error {
	req := portalapi.DeletePoolRequest{
		PoolName: poolName,
	}

	hReq, err := c.createPostRequestWithJsonBody(context.Background(), "/api/v1/delete-pool", nil, req)
	if err != nil {
		return fmt.Errorf("failed to delete pool: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to delete pool: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to delete pool: %w", err)
	}

	return nil
}
