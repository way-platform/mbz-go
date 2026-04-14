package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
	fleetv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/fleet/v1"
)

// GetVehicleCompatibility gets the data service compatibility for a specific vehicle.
func (c *Client) GetVehicleCompatibility(
	ctx context.Context,
	request *fleetv1.GetVehicleCompatibilityRequest,
) (_ *fleetv1.GetVehicleCompatibilityResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle compatibility: %w", err)
		}
	}()
	requestURL, err := url.JoinPath(
		c.baseURL,
		"/v1/accounts/vehicles",
		request.GetVin(),
		"compatibilities",
	)
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
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
	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var responseBody vehiclesv1.CompatibilityResponse
	if err := json.Unmarshal(data, &responseBody); err != nil {
		return nil, err
	}
	result := compatibilityResponseToProto(request.GetVin(), &responseBody)
	if rawStruct, err := compatibilityRawStruct(data); err == nil {
		result.SetRaw(rawStruct)
	}
	resp := &fleetv1.GetVehicleCompatibilityResponse{}
	resp.SetVehicleCompatibility(result)
	return resp, nil
}
