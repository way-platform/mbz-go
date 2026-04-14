package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	vehiclespecv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/vehiclespec/v1"
)

// GetVehicleSpecification gets the vehicle marketing information for a given VIN.
// This method requires API key authentication via [WithAPIKey].
func (c *Client) GetVehicleSpecification(
	ctx context.Context,
	request *vehiclespecv1.GetVehicleSpecificationRequest,
) (_ *vehiclespecv1.GetVehicleSpecificationResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle specification: %w", err)
		}
	}()
	if request.GetVin() == "" {
		return nil, fmt.Errorf("VIN is required")
	}
	values := url.Values{}
	if request.GetLocale() != "" {
		values.Set("locale", request.GetLocale())
	} else {
		values.Set("locale", string(vehiclespecificationfleetv1.LocalesEnUS))
	}
	// Set all optional parameters to true to maximize data retrieval.
	values.Set("additionalSpecs", "true")
	values.Set("optionsNullDescription", "true")
	values.Set("options", "true")
	values.Set("technicalData", "true")
	values.Set("payloadNullValues", "true")
	requestURL := fmt.Sprintf("%s/vehicles/%s", vehiclespecificationfleetv1.BaseURL, request.GetVin())
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.URL.RawQuery = values.Encode()
	httpRequest.Header.Set("Accept", "application/json")
	httpResponse, err := c.httpClient(c.config).Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := httpResponse.Body.Close(); closeErr != nil {
			log.Printf("mbz: failed to close response body: %v", closeErr)
		}
	}()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	responseData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var openAPIResp vehiclespecificationfleetv1.VehicleSpecificationResponse
	if err := json.Unmarshal(responseData, &openAPIResp); err != nil {
		return nil, err
	}
	spec := vehicleDataToProto(openAPIResp.VehicleData)
	if rawStruct, err := vehicleDataRawStruct(responseData); err == nil {
		spec.SetRaw(rawStruct)
	}
	resp := &vehiclespecv1.GetVehicleSpecificationResponse{}
	resp.SetVehicleSpecification(spec)
	return resp, nil
}
