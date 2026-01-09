package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nvvers/orbit-portal-client-golang/internal/orbit"
	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
)

func (c *Client) PullNotification(poolName string) (*orbit.Notification, error) {
	req := portalapi.PullNotificationRequest{PoolName: poolName}
	hReq, err := c.createPostRequest(context.Background(), "/api/v1/pull-notification", "", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for pulling notification: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return nil, fmt.Errorf("failed to pull notification: %w", err)
	}
	defer hRes.Body.Close()

	if hRes.StatusCode == http.StatusOK {
		var n orbit.Notification
		if err := json.NewDecoder(hRes.Body).Decode(&n); err != nil {
			return nil, fmt.Errorf("failed to unmarshal n: %w", err)
		}

		return &n, nil

	} else if hRes.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	resData, _ := io.ReadAll(hRes.Body)
	return nil, fmt.Errorf("failed to pull notification, code: %d: %s", hRes.StatusCode, string(resData))
}
