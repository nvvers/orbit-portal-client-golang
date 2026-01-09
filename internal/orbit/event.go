package orbit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event repräsentiert ein Ereignis
type Event struct {
	Source      string       `json:"source"`
	Subject     string       `json:"subject"`
	Type        string       `json:"type"`
	Data        any          `json:"data"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// EventRecord repräsentiert ein Ereignis, das von einem Event-Pool verarbeitet und als Notification an die verbundenen Clients gesendet wird.
type EventRecord struct {
	ID          uuid.UUID       `json:"id"`
	Time        time.Time       `json:"time"`
	Source      string          `json:"source"`
	Subject     string          `json:"subject"`
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
	Attachments []Attachment    `json:"attachments,omitempty"`
}
