package orbit

import (
	"time"
)

// Notification repräsentiert eine Benachrichtigung, die von einem Event-Pool an verbundene Clients gesendet wird.
type Notification struct {
	// Der Name des Pools.
	PoolName string `json:"poolName"`

	// Das Ereignis, das die Benachrichtigung ausgelöst hat.
	Event EventRecord `json:"event"`

	// Gibt an, wann eine Benachrichtigung verfällt.
	Expiration time.Time `json:"expiration,omitempty"`

	// Gibt an, wann die nächste Zustellung versucht wird.
	NextDeliveryTime time.Time `json:"nextDeliveryTime,omitempty"`

	// Zählt, wie oft die Benachrichtigung zugestellt wurde.
	DeliveryCount int `json:"deliveryCount"`
}
