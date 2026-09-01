package portalapi

import (
	"github.com/nvvers/orbit-portal-client-golang/orbit"

	"github.com/google/uuid"
)

type (
	DiscardNotificationRequest struct {
		PoolName         string        `json:"poolName"`
		EventID          uuid.UUID     `json:"eventId"`
		SubsequentEvents []orbit.Event `json:"subsequentEvents,omitempty"`
	}

	DiscardNotificationResponse struct{}
)
