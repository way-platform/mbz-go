package mbz

import (
	"encoding/json"
	"fmt"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ParseRawVehicleSpecification reconstructs a [mbzv1.VehicleSpecification]
// from a raw [structpb.Struct] previously captured by [Client.GetVehicleSpecification].
// This enables round-tripping: store the raw struct, then re-parse it later
// through the same conversion pipeline used by the live API call.
func ParseRawVehicleSpecification(raw *structpb.Struct) (*mbzv1.VehicleSpecification, error) {
	if raw == nil {
		return nil, fmt.Errorf("mbz: parse vehicle specification from raw: nil struct")
	}
	jsonBytes, err := json.Marshal(raw.AsMap())
	if err != nil {
		return nil, fmt.Errorf("mbz: parse vehicle specification from raw: marshal: %w", err)
	}
	var vd vehiclespecificationfleetv1.VehicleData
	if err := json.Unmarshal(jsonBytes, &vd); err != nil {
		return nil, fmt.Errorf("mbz: parse vehicle specification from raw: unmarshal: %w", err)
	}
	return vehicleDataToProto(&vd), nil
}

// vehicleDataRawStruct extracts the vehicleData portion of the API response
// as a [structpb.Struct] for storage and later re-parsing.
func vehicleDataRawStruct(responseBody []byte) (*structpb.Struct, error) {
	var envelope struct {
		VehicleData json.RawMessage `json:"vehicleData"`
	}
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return nil, err
	}
	if len(envelope.VehicleData) == 0 {
		return nil, fmt.Errorf("no vehicleData in response")
	}
	var m map[string]any
	if err := json.Unmarshal(envelope.VehicleData, &m); err != nil {
		return nil, err
	}
	return structpb.NewStruct(m)
}
