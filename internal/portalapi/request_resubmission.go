package portalapi

import (
	"time"

	"github.com/google/uuid"
)

type (
	RequestResubmissionRequest struct {
		PoolName         string    `json:"poolName"`
		EventID          uuid.UUID `json:"eventId"`
		ResubmissionTime time.Time `json:"resubmissionTime"`
	}

	RequestResubmissionResponse struct{}
)
