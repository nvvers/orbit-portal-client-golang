package orbit

import "encoding/json"

type MessageType string

const (
	MessageTypeHeartbeat    MessageType = "heartbeat"
	MessageTypeNotification MessageType = "notification"
	MessageTypeClose        MessageType = "close"
)

// Message repräsentiert eine Nachricht, die vom Portal-Server an einen observierenden Client gesendet wird.
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
