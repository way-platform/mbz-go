package mbz

import (
	"encoding/json"

	"github.com/way-platform/mbz-go/api/servicesv1"
)

// Message represents a push message.
//
// See: https://developer.mercedes-benz.com/products/connect_your_fleet/specifications/push_api
type Message struct {
	// Unique message identifier.
	MessageID string `json:"messageId"`

	// Vehicle identification number.
	VIN string `json:"vin"`

	// Time when the message was created (in milliseconds since 1970).
	Timestamp int64 `json:"timestamp"`

	// Always vehiclesignal.
	MessageType string `json:"messageType"`

	// Version tag for the message.
	Version string `json:"version"`

	// Service associated with the message.
	ServiceID string `json:"serviceId"`

	// Sending mode.
	SendingBehaviour servicesv1.SignalSendingBehaviour `json:"sendingBehaviour"`

	// Message data.
	Data json.RawMessage `json:"data"`
}
