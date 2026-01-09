package portalapi

type (
	SaveSchemaRequest struct {
		EventType string         `json:"eventType"`
		Schema    map[string]any `json:"schema"`
	}

	SaveSchemaResponse struct{}
)
