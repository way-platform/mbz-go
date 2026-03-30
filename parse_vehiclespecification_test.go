package mbz

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestParseRawVehicleSpecification_roundTrip(t *testing.T) {
	testdataDir := "testdata/specifications"
	err := filepath.WalkDir(testdataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.HasSuffix(path, ".golden") || filepath.Ext(path) != ".json" {
			return nil
		}
		relPath, err := filepath.Rel(testdataDir, path)
		if err != nil {
			return err
		}
		t.Run(relPath, func(t *testing.T) {
			inputData, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read input: %v", err)
			}

			// Normal path: JSON -> OpenAPI -> proto
			var openAPIResp vehiclespecificationfleetv1.VehicleSpecificationResponse
			if err := json.Unmarshal(inputData, &openAPIResp); err != nil {
				t.Fatalf("unmarshal input: %v", err)
			}
			original := vehicleDataToProto(openAPIResp.VehicleData)

			// Extract raw struct from the response body (same as GetVehicleSpecification does)
			rawStruct, err := vehicleDataRawStruct(inputData)
			if err != nil {
				t.Fatalf("vehicleDataRawStruct: %v", err)
			}

			// Round-trip: raw struct -> parsed proto
			roundTripped, err := ParseRawVehicleSpecification(rawStruct)
			if err != nil {
				t.Fatalf("ParseRawVehicleSpecification: %v", err)
			}

			// Clear the raw field on original before comparison (round-tripped won't have it)
			original.ClearRaw()
			roundTripped.ClearRaw()

			if !proto.Equal(original, roundTripped) {
				origJSON, _ := protojson.Marshal(original)
				rtJSON, _ := protojson.Marshal(roundTripped)
				t.Errorf("round-trip mismatch\noriginal:     %s\nround-tripped: %s", origJSON, rtJSON)
			}
		})
		return nil
	})
	if err != nil {
		t.Fatalf("walk testdata: %v", err)
	}
}
