package portalapi

import "github.com/nvvers/orbit-portal-client-golang/internal/orbit"

type (
	PullNotificationRequest struct {
		PoolName string `json:"poolName"`
	}

	PullNotificationResponse orbit.Notification
)
