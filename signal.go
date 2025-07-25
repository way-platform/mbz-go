package mbz

import (
	"fmt"
	"strconv"
	"strings"

	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/mbz/v1"
	"google.golang.org/protobuf/proto"
)

type VehicleSignalData struct {
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

func signalNameToIdentifierEnum(name string) (mbzv1.SignalIdentifier, error) {
	id := signalNameToIdentifier(name)
	enumValue, ok := mbzv1.SignalIdentifier_value[id]
	if !ok {
		return mbzv1.SignalIdentifier_SIGNAL_IDENTIFIER_UNSPECIFIED, fmt.Errorf("unknown signal name: %s (id: %s)", name, id)
	}
	return mbzv1.SignalIdentifier(enumValue), nil
}

func unixTimestampMillisToMicros(timestamp int64) int64 {
	return timestamp * 1000
}

func (s *Signal) AsProto() (*mbzv1.Signal, error) {
	identifier, err := signalNameToIdentifierEnum(s.Name)
	if err != nil {
		return nil, err
	}
	result := &mbzv1.Signal{
		Id:   identifier,
		Time: unixTimestampMillisToMicros(s.Timestamp),
	}
	signalType, ok := proto.GetExtension(
		identifier.Descriptor().Values().ByNumber(identifier.Number()).Options(),
		mbzv1.E_SignalType,
	).(mbzv1.SignalType)
	if !ok {
		return nil, fmt.Errorf("failed to get signal type for signal %s", s.Name)
	}
	switch signalType {
	case mbzv1.SignalType_STRING:
		result.StringValue = ptr(s.Value)
	case mbzv1.SignalType_INTEGER:
		intValue, err := strconv.Atoi(s.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse integer value for signal %s: %w", s.Name, err)
		}
		result.IntValue = ptr(int32(intValue))
	case mbzv1.SignalType_DOUBLE:
		doubleValue, err := strconv.ParseFloat(s.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse double value for signal %s: %w", s.Name, err)
		}
		result.DoubleValue = ptr(doubleValue)
	case mbzv1.SignalType_BOOLEAN:
		booleanValue, err := strconv.ParseBool(s.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse boolean value for signal %s: %w", s.Name, err)
		}
		result.BoolValue = ptr(booleanValue)
	case mbzv1.SignalType_ENUM:
		enumValues, ok := proto.GetExtension(
			identifier.Descriptor().Values().ByNumber(identifier.Number()).Options(),
			mbzv1.E_SignalValues,
		).([]*mbzv1.SignalEnumValue)
		if !ok {
			return nil, fmt.Errorf("failed to get enum values for signal %s", s.Name)
		}
		var validEnumValue bool
		for _, enumValue := range enumValues {
			if enumValue.GetValue() == s.Value {
				validEnumValue = true
				break
			}
		}
		if !validEnumValue {
			return nil, fmt.Errorf("invalid enum value for signal %s: %s", s.Name, s.Value)
		}
		result.EnumValue = ptr(s.Value)
	default:
		return nil, fmt.Errorf("unknown signal type: %s", signalType)
	}
	signalUnit, ok := proto.GetExtension(
		identifier.Descriptor().Values().ByNumber(identifier.Number()).Options(),
		mbzv1.E_SignalUnit,
	).(mbzv1.SignalUnit)
	if !ok {
		return nil, fmt.Errorf("failed to get signal unit for signal %s", s.Name)
	}
	if signalUnit != mbzv1.SignalUnit_SIGNAL_UNIT_UNSPECIFIED {
		result.Unit = ptr(signalUnit)
	}
	return result, nil
}

func ptr[T any](v T) *T {
	return &v
}
