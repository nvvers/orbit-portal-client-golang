package portal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nvvers/orbit-portal-client-golang/internal/portalapi"
)

type PoolOption func(*portalapi.SavePoolRequest) error

func (c *Client) SavePool(poolName string, opts ...PoolOption) error {
	req := portalapi.SavePoolRequest{
		PoolName:       poolName,
		SubjectPattern: ".*",
		// Set default values
		EventTypes:                []string{},
		NotificationRetryInterval: int((5 * time.Minute).Seconds()),
		NotificationTTL:           int((24 * time.Hour).Seconds()),
	}

	for _, opt := range opts {
		if err := opt(&req); err != nil {
			return err
		}
	}

	hReq, err := c.createPostRequest(context.Background(), "/api/v1/save-pool", "", req)
	if err != nil {
		return fmt.Errorf("failed to save pool: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to save pool: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to save pool: %w", err)
	}

	return nil
}

func WithSubjectPattern(pattern string) PoolOption {
	return func(req *portalapi.SavePoolRequest) error {
		req.SubjectPattern = pattern
		return nil
	}
}

func WithEventTypes(eventTypes []string) PoolOption {
	return func(req *portalapi.SavePoolRequest) error {
		req.EventTypes = eventTypes
		return nil
	}
}

func WithNotificationRetryInterval(interval time.Duration) PoolOption {
	return func(req *portalapi.SavePoolRequest) error {
		if interval <= 0 {
			return fmt.Errorf("notification retry interval must be greater than 0")
		}

		req.NotificationRetryInterval = int(interval.Seconds())
		return nil
	}
}

func WithNotificationTTL(ttl time.Duration) PoolOption {
	return func(req *portalapi.SavePoolRequest) error {
		if ttl <= 0 {
			return fmt.Errorf("notification ttl must be greater than 0")
		}

		req.NotificationTTL = int(ttl.Seconds())
		return nil
	}
}
