package mbz

import (
	"encoding/json"
	"fmt"

	"github.com/way-platform/mbz-go/api/servicesv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/mbz/v1"
)

// PushMessage represents a push message.
//
// See: https://developer.mercedes-benz.com/products/connect_your_fleet/specifications/push_api
type PushMessage struct {
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

// AsProto converts the push message to an [mbzv1.PushMessage] protobuf message.
func (m *PushMessage) AsProto() (*mbzv1.PushMessage, error) {
	messageType, err := messageTypeToProto(m.MessageType)
	if err != nil {
		return nil, err
	}
	sendingBehavior, err := sendingBehaviorToProto(m.SendingBehaviour)
	if err != nil {
		return nil, err
	}
	var result mbzv1.PushMessage
	result.SetMessageId(m.MessageID)
	result.SetVin(m.VIN)
	result.SetTime(unixTimestampMillisToMicros(m.Timestamp))
	result.SetMessageType(messageType)
	result.SetVersion(m.Version)
	result.SetServiceId(m.ServiceID)
	result.SetSendingBehavior(sendingBehavior)
	switch messageType {
	case mbzv1.MessageType_SIGNALS:
		var data VehicleSignalData
		if err := json.Unmarshal(m.Data, &data); err != nil {
			return nil, err
		}
		result.SetSignals(make([]*mbzv1.Signal, 0, len(data.Signals)))
		for _, signal := range data.Signals {
			signalProto, err := signal.AsProto()
			if err != nil {
				return nil, err
			}
			result.SetSignals(append(result.GetSignals(), signalProto))
		}
	}
	return &result, nil
}

func sendingBehaviorToProto(sendingBehavior servicesv1.SignalSendingBehaviour) (mbzv1.SendingBehavior, error) {
	switch sendingBehavior {
	case servicesv1.SignalSendingBehaviourONCHANGE:
		return mbzv1.SendingBehavior_ON_CHANGE, nil
	case servicesv1.SignalSendingBehaviourONCHARGINGCYCLE:
		return mbzv1.SendingBehavior_ON_CHARGING_CYCLE, nil
	case servicesv1.SignalSendingBehaviourONCHARGINGCYCLEEND:
		return mbzv1.SendingBehavior_ON_CHARGING_CYCLE_END, nil
	case servicesv1.SignalSendingBehaviourONCHARGINGDETECTION:
		return mbzv1.SendingBehavior_ON_CHARGING_DETECTION, nil
	case servicesv1.SignalSendingBehaviourONCHARGINGSESSION:
		return mbzv1.SendingBehavior_ON_CHARGING_SESSION, nil
	case servicesv1.SignalSendingBehaviourONCHARGINGSESSIONEND:
		return mbzv1.SendingBehavior_ON_CHARGING_SESSION_END, nil
	case servicesv1.SignalSendingBehaviourONIGNITIONLOCK:
		return mbzv1.SendingBehavior_ON_IGNITION_LOCK, nil
	case servicesv1.SignalSendingBehaviourONINTERVAL120SEC:
		return mbzv1.SendingBehavior_ON_INTERVAL_120_SEC, nil
	case servicesv1.SignalSendingBehaviourONINTERVAL15SEC:
		return mbzv1.SendingBehavior_ON_INTERVAL_15_SEC, nil
	case servicesv1.SignalSendingBehaviourONINTERVAL30SEC:
		return mbzv1.SendingBehavior_ON_INTERVAL_30_SEC, nil
	case servicesv1.SignalSendingBehaviourONRECHARGESESSION:
		return mbzv1.SendingBehavior_ON_RECHARGE_SESSION, nil
	case servicesv1.SignalSendingBehaviourONREFUELINGEND:
		return mbzv1.SendingBehavior_ON_REFUELING_END, nil
	case servicesv1.SignalSendingBehaviourONREFUELSESSION:
		return mbzv1.SendingBehavior_ON_REFUEL_SESSION, nil
	case servicesv1.SignalSendingBehaviourONTRIP:
		return mbzv1.SendingBehavior_ON_TRIP, nil
	case servicesv1.SignalSendingBehaviourONTRIPEND:
		return mbzv1.SendingBehavior_ON_TRIP_END, nil
	case servicesv1.SignalSendingBehaviourONREFRESH:
		return mbzv1.SendingBehavior_ON_REFRESH, nil
	default:
		return mbzv1.SendingBehavior_SENDING_BEHAVIOR_UNSPECIFIED, fmt.Errorf("unknown sending behavior: %s", sendingBehavior)
	}
}

func messageTypeToProto(messageType string) (mbzv1.MessageType, error) {
	switch messageType {
	case "vehiclesignal":
		return mbzv1.MessageType_SIGNALS, nil
	case "vehicle_registration_response":
		return mbzv1.MessageType_REGISTRATION_RESPONSE, nil
	case "vehicle_service_status_update":
		return mbzv1.MessageType_SERVICE_STATUS_UPDATE, nil
	case "vehicle_service_activation_pending":
		return mbzv1.MessageType_SERVICE_ACTIVATION_PENDING, nil
	case "vehicle_trip_summary":
		return mbzv1.MessageType_TRIP_SUMMARY, nil
	case "refueling_detection":
		return mbzv1.MessageType_REFUELING_DETECTION, nil
	case "charging_detection":
		return mbzv1.MessageType_CHARGING_DETECTION, nil
	case "charging_cycle_summary":
		return mbzv1.MessageType_CHARGING_CYCLE_SUMMARY, nil
	default:
		return mbzv1.MessageType_MESSAGE_TYPE_UNSPECIFIED, fmt.Errorf("unknown message type: %s", messageType)
	}
}
