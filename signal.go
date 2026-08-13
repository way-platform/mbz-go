package mbz

import (
	"strconv"
	"strings"
	"time"

	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// VehicleSignalData is the decoded payload of a vehicle signal push message.
type VehicleSignalData struct {
	// Signals contained in the message.
	Signals []Signal `json:"signals"`
}

// Signal represents a vehicle signal.
type Signal struct {
	// Name of the signal.
	Name string `json:"name"`

	// Timestamp of the signal (unix timestamp in milliseconds).
	Timestamp int64 `json:"timestamp"`

	// Value of the signal.
	Value string `json:"value"`

	// Type of the signal.
	Type string `json:"type"`
}

func signalNameToIdentifier(name string) string {
	return strings.ReplaceAll(strings.ToUpper(name), ".", "_")
}

func signalNameToIdentifierEnum(name string) (mbzv1.SignalIdentifier, bool) {
	id := signalNameToIdentifier(name)
	enumValue, ok := mbzv1.SignalIdentifier_value[id]
	if !ok {
		return mbzv1.SignalIdentifier_SIGNAL_IDENTIFIER_UNSPECIFIED, false
	}
	return mbzv1.SignalIdentifier(enumValue), true
}

func (s *Signal) unmapped(reason mbzv1.UnmappedSignal_Reason) *mbzv1.UnmappedSignal {
	var result mbzv1.UnmappedSignal
	result.SetName(s.Name)
	result.SetTime(timestamppb.New(time.UnixMilli(s.Timestamp).UTC()))
	result.SetValue(s.Value)
	result.SetType(s.Type)
	result.SetReason(reason)
	return &result
}

func (s *Signal) convertToProto() (*mbzv1.Signal, *mbzv1.UnmappedSignal) {
	identifier, ok := signalNameToIdentifierEnum(s.Name)
	if !ok {
		return nil, s.unmapped(mbzv1.UnmappedSignal_UNKNOWN_NAME)
	}
	var result mbzv1.Signal
	result.SetId(identifier)
	result.SetTime(timestamppb.New(time.UnixMilli(s.Timestamp).UTC()))
	signalType, ok := proto.GetExtension(
		identifier.Descriptor().Values().ByNumber(identifier.Number()).Options(),
		mbzv1.E_SignalType,
	).(mbzv1.SignalType)
	if !ok {
		return nil, s.unmapped(mbzv1.UnmappedSignal_PARSE_ERROR)
	}
	result.SetType(signalType)
	switch signalType {
	case mbzv1.SignalType_STRING:
		result.SetStringValue(s.Value)
	case mbzv1.SignalType_INTEGER:
		intValue, err := strconv.ParseInt(s.Value, 10, 32)
		if err != nil {
			return nil, s.unmapped(mbzv1.UnmappedSignal_PARSE_ERROR)
		}
		result.SetIntValue(int32(intValue))
	case mbzv1.SignalType_DOUBLE:
		doubleValue, err := strconv.ParseFloat(s.Value, 64)
		if err != nil {
			return nil, s.unmapped(mbzv1.UnmappedSignal_PARSE_ERROR)
		}
		result.SetDoubleValue(doubleValue)
	case mbzv1.SignalType_BOOLEAN:
		booleanValue, err := strconv.ParseBool(s.Value)
		if err != nil {
			return nil, s.unmapped(mbzv1.UnmappedSignal_PARSE_ERROR)
		}
		result.SetBoolValue(booleanValue)
	case mbzv1.SignalType_ENUM:
		enumValues, ok := proto.GetExtension(
			identifier.Descriptor().Values().ByNumber(identifier.Number()).Options(),
			mbzv1.E_SignalValues,
		).([]*mbzv1.SignalEnumValue)
		if !ok {
			return nil, s.unmapped(mbzv1.UnmappedSignal_PARSE_ERROR)
		}
		var validEnumValue bool
		for _, enumValue := range enumValues {
			if enumValue.GetValue() == s.Value {
				validEnumValue = true
				break
			}
		}
		if !validEnumValue {
			return nil, s.unmapped(mbzv1.UnmappedSignal_INVALID_ENUM)
		}
		result.SetEnumValue(s.Value)
	default:
		return nil, s.unmapped(mbzv1.UnmappedSignal_PARSE_ERROR)
	}
	signalUnit, ok := proto.GetExtension(
		identifier.Descriptor().Values().ByNumber(identifier.Number()).Options(),
		mbzv1.E_SignalUnit,
	).(mbzv1.SignalUnit)
	if !ok {
		return nil, s.unmapped(mbzv1.UnmappedSignal_PARSE_ERROR)
	}
	if signalUnit != mbzv1.SignalUnit_SIGNAL_UNIT_UNSPECIFIED {
		result.SetUnit(signalUnit)
	}
	return &result, nil
}
