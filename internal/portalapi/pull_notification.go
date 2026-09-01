package portalapi

import (
	"github.com/nvvers/orbit-portal-client-golang/orbit"
)

type (
	PullNotificationRequest struct {
		PoolName string `json:"poolName"`
	}

	PullNotificationResponse orbit.Notification
)
