package mbz

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	"google.golang.org/protobuf/encoding/protojson"
)

var update = flag.Bool("update", false, "update golden files")

func Test_vehicleDataToProto_golden(t *testing.T) {
	testdataDir := "testdata/specifications"
	err := filepath.WalkDir(testdataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// Skip golden files
		if strings.HasSuffix(path, ".golden") {
			return nil
		}
		// Only process JSON files
		if filepath.Ext(path) != ".json" {
			return nil
		}
		relPath, err := filepath.Rel(testdataDir, path)
		if err != nil {
			return err
		}
		t.Run(relPath, func(t *testing.T) {
			goldenPath := path + ".golden"
			// Read input JSON
			inputData, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read input file: %v", err)
			}
			// Unmarshal into OpenAPI response
			var openAPIResp vehiclespecificationfleetv1.VehicleSpecificationResponse
			if err := json.Unmarshal(inputData, &openAPIResp); err != nil {
				t.Fatalf("failed to unmarshal input JSON: %v", err)
			}
			// Convert to proto
			protoSpec := vehicleDataToProto(openAPIResp.VehicleData)
			// Marshal proto to JSON
			actualJSONBytes, err := protojson.Marshal(protoSpec)
			if err != nil {
				t.Fatalf("failed to marshal proto to JSON: %v", err)
			}
			// Format JSON consistently using json.Indent
			var actualJSONBuf bytes.Buffer
			if err := json.Indent(&actualJSONBuf, actualJSONBytes, "", "  "); err != nil {
				t.Fatalf("failed to indent JSON: %v", err)
			}
			actualJSON := actualJSONBuf.Bytes()
			// Update golden file if flag is set
			if *update {
				if err := os.WriteFile(goldenPath, actualJSON, 0o644); err != nil {
					t.Fatalf("failed to write golden file: %v", err)
				}
				t.Logf("updated golden file: %s", goldenPath)
				return
			}
			// Read golden file
			expectedJSON, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden file: %v (run with -update to create)", err)
			}
			// Compare (normalize JSON for comparison)
			var actual, expected interface{}
			if err := json.Unmarshal(actualJSON, &actual); err != nil {
				t.Fatalf("failed to unmarshal actual JSON: %v", err)
			}
			if err := json.Unmarshal(expectedJSON, &expected); err != nil {
				t.Fatalf("failed to unmarshal expected JSON: %v", err)
			}
			// Re-marshal for comparison (normalizes formatting)
			actualNormalized, err := json.MarshalIndent(actual, "", "  ")
			if err != nil {
				t.Fatalf("failed to normalize actual JSON: %v", err)
			}
			expectedNormalized, err := json.MarshalIndent(expected, "", "  ")
			if err != nil {
				t.Fatalf("failed to normalize expected JSON: %v", err)
			}

			if string(actualNormalized) != string(expectedNormalized) {
				t.Errorf("conversion result differs from golden file\n\nActual:\n%s\n\nExpected:\n%s", string(actualNormalized), string(expectedNormalized))
			}
		})
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk testdata directory: %v", err)
	}
}
