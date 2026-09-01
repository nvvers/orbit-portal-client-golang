package portalapi

import (
	"github.com/nvvers/orbit-portal-client-golang/orbit"

	"github.com/google/uuid"
)

type (
	AcknowledgeNotificationRequest struct {
		PoolName         string        `json:"poolName"`
		EventID          uuid.UUID     `json:"eventId"`
		SubsequentEvents []orbit.Event `json:"subsequentEvents,omitempty"`
	}

	AcknowledgeNotificationResponse struct{}
)
