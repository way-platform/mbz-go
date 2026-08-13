package mbz

import (
	"encoding/json"
	"time"

	"github.com/way-platform/mbz-go/api/servicesv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
//
// Returns an error only when the signal data JSON cannot be unmarshaled.
// Unknown signal names, invalid enum values, and unparseable values are
// collected in [mbzv1.PushMessage.unmapped_signals].
func (m *PushMessage) AsProto() (*mbzv1.PushMessage, error) {
	messageType := messageTypeToProto(m.MessageType)
	var result mbzv1.PushMessage
	result.SetMessageId(m.MessageID)
	result.SetMessageType(messageType)
	result.SetTime(timestamppb.New(time.UnixMilli(m.Timestamp).UTC()))
	result.SetVersion(m.Version)
	if messageType != mbzv1.MessageType_SIGNALS {
		return &result, nil
	}
	result.SetVin(m.VIN)
	result.SetServiceId(m.ServiceID)
	result.SetSendingBehavior(sendingBehaviorToProto(m.SendingBehaviour))
	var data VehicleSignalData
	if err := json.Unmarshal(m.Data, &data); err != nil {
		return nil, err
	}
	signals := make([]*mbzv1.Signal, 0, len(data.Signals))
	unmappedSignals := make([]*mbzv1.UnmappedSignal, 0)
	for _, signal := range data.Signals {
		signalProto, unmappedProto := signal.convertToProto()
		if signalProto != nil {
			signals = append(signals, signalProto)
		}
		if unmappedProto != nil {
			unmappedSignals = append(unmappedSignals, unmappedProto)
		}
	}
	result.SetSignals(signals)
	result.SetUnmappedSignals(unmappedSignals)
	return &result, nil
}

func sendingBehaviorToProto(
	sendingBehavior servicesv1.SignalSendingBehaviour,
) mbzv1.SendingBehavior {
	switch sendingBehavior {
	case servicesv1.SignalSendingBehaviourONCHANGE:
		return mbzv1.SendingBehavior_ON_CHANGE
	case servicesv1.SignalSendingBehaviourONCHARGINGCYCLE:
		return mbzv1.SendingBehavior_ON_CHARGING_CYCLE
	case servicesv1.SignalSendingBehaviourONCHARGINGCYCLEEND:
		return mbzv1.SendingBehavior_ON_CHARGING_CYCLE_END
	case servicesv1.SignalSendingBehaviourONCHARGINGDETECTION:
		return mbzv1.SendingBehavior_ON_CHARGING_DETECTION
	case servicesv1.SignalSendingBehaviourONCHARGINGSESSION:
		return mbzv1.SendingBehavior_ON_CHARGING_SESSION
	case servicesv1.SignalSendingBehaviourONCHARGINGSESSIONEND:
		return mbzv1.SendingBehavior_ON_CHARGING_SESSION_END
	case servicesv1.SignalSendingBehaviourONIGNITIONLOCK:
		return mbzv1.SendingBehavior_ON_IGNITION_LOCK
	case servicesv1.SignalSendingBehaviourONINTERVAL120SEC:
		return mbzv1.SendingBehavior_ON_INTERVAL_120_SEC
	case servicesv1.SignalSendingBehaviourONINTERVAL15SEC:
		return mbzv1.SendingBehavior_ON_INTERVAL_15_SEC
	case servicesv1.SignalSendingBehaviourONINTERVAL30SEC:
		return mbzv1.SendingBehavior_ON_INTERVAL_30_SEC
	case servicesv1.SignalSendingBehaviourONRECHARGESESSION:
		return mbzv1.SendingBehavior_ON_RECHARGE_SESSION
	case servicesv1.SignalSendingBehaviourONREFUELINGEND:
		return mbzv1.SendingBehavior_ON_REFUELING_END
	case servicesv1.SignalSendingBehaviourONREFUELSESSION:
		return mbzv1.SendingBehavior_ON_REFUEL_SESSION
	case servicesv1.SignalSendingBehaviourONTRIP:
		return mbzv1.SendingBehavior_ON_TRIP
	case servicesv1.SignalSendingBehaviourONTRIPEND:
		return mbzv1.SendingBehavior_ON_TRIP_END
	case servicesv1.SignalSendingBehaviourONREFRESH:
		return mbzv1.SendingBehavior_ON_REFRESH
	default:
		return mbzv1.SendingBehavior_SENDING_BEHAVIOR_UNSPECIFIED
	}
}

func messageTypeToProto(messageType string) mbzv1.MessageType {
	switch messageType {
	case "vehiclesignal":
		return mbzv1.MessageType_SIGNALS
	case "vehicletrip":
		return mbzv1.MessageType_TRIP_SUMMARY
	case "vehicle_registration_response":
		return mbzv1.MessageType_REGISTRATION_RESPONSE
	case "vehicle_service_status_update":
		return mbzv1.MessageType_SERVICE_STATUS_UPDATE
	case "vehicle_service_activation_pending":
		return mbzv1.MessageType_SERVICE_ACTIVATION_PENDING
	case "refueling_detection":
		return mbzv1.MessageType_REFUELING_DETECTION
	case "charging_detection":
		return mbzv1.MessageType_CHARGING_DETECTION
	case "charging_cycle_summary":
		return mbzv1.MessageType_CHARGING_CYCLE_SUMMARY
	case "vehicle_command_response":
		return mbzv1.MessageType_COMMAND_RESPONSE
	default:
		return mbzv1.MessageType_MESSAGE_TYPE_UNSPECIFIED
	}
}
