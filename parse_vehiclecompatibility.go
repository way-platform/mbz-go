package mbz

import (
	"encoding/json"
	"fmt"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ParseRawVehicleCompatibility reconstructs a [mbzv1.VehicleCompatibility]
// from a raw [structpb.Struct] previously captured by [Client.GetVehicleCompatibility].
func ParseRawVehicleCompatibility(raw *structpb.Struct) (*mbzv1.VehicleCompatibility, error) {
	if raw == nil {
		return nil, fmt.Errorf("mbz: parse vehicle compatibility from raw: nil struct")
	}
	jsonBytes, err := json.Marshal(raw.AsMap())
	if err != nil {
		return nil, fmt.Errorf("mbz: parse vehicle compatibility from raw: marshal: %w", err)
	}
	var resp vehiclesv1.CompatibilityResponse
	if err := json.Unmarshal(jsonBytes, &resp); err != nil {
		return nil, fmt.Errorf("mbz: parse vehicle compatibility from raw: unmarshal: %w", err)
	}
	// VIN is not present in the response body — pass empty, caller knows the VIN.
	return compatibilityResponseToProto("", &resp), nil
}

// compatibilityRawStruct converts a compatibility API response body into a [structpb.Struct].
func compatibilityRawStruct(responseBody []byte) (*structpb.Struct, error) {
	var m map[string]any
	if err := json.Unmarshal(responseBody, &m); err != nil {
		return nil, err
	}
	return structpb.NewStruct(m)
}
