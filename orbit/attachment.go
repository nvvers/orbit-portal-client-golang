package orbit

import "github.com/google/uuid"

// Attachment repräsentiert eine Datei, die an ein Ereignis angehängt werden kann.
type Attachment struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Description string    `json:"description"`
}
