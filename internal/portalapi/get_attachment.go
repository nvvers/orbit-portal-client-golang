package portalapi

import (
	"github.com/google/uuid"
)

type (
	GetAttachmentRequest struct {
		AttachmentID uuid.UUID `json:"attachmentId"`
	}
)
