package mbz

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestParseRawVehicleCompatibility_roundTrip(t *testing.T) {
	testdataDir := "testdata/compatibilities"
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
			var resp vehiclesv1.CompatibilityResponse
			if err := json.Unmarshal(inputData, &resp); err != nil {
				t.Fatalf("unmarshal input: %v", err)
			}
			original := compatibilityResponseToProto(resp.VIN, &resp)

			// Extract raw struct from the response body
			rawStruct, err := compatibilityRawStruct(inputData)
			if err != nil {
				t.Fatalf("compatibilityRawStruct: %v", err)
			}

			// Round-trip: raw struct -> parsed proto
			roundTripped, err := ParseRawVehicleCompatibility(rawStruct)
			if err != nil {
				t.Fatalf("ParseRawVehicleCompatibility: %v", err)
			}

			// Clear the raw field on original before comparison (round-tripped won't have it).
			// Also set VIN on round-tripped since ParseRaw passes empty VIN.
			original.ClearRaw()
			roundTripped.ClearRaw()
			roundTripped.SetVin(original.GetVin())

			if !proto.Equal(original, roundTripped) {
				origJSON, _ := protojson.Marshal(original)
				rtJSON, _ := protojson.Marshal(roundTripped)
				t.Errorf(
					"round-trip mismatch\noriginal:     %s\nround-tripped: %s",
					origJSON,
					rtJSON,
				)
			}
		})
		return nil
	})
	if err != nil {
		t.Fatalf("walk testdata: %v", err)
	}
}
