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
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
)

// GetVehicleCompatibilityRequest is the request for [Client.GetVehicleCompatibility].
type GetVehicleCompatibilityRequest struct {
	// VIN of the vehicle to get the compatibility for.
	VIN string `json:"vin"`
}

// GetVehicleCompatibility gets the compatibility of a vehicle.
func (c *Client) GetVehicleCompatibility(
	ctx context.Context,
	request *GetVehicleCompatibilityRequest,
	opts ...ClientOption,
) (_ *mbzv1.VehicleCompatibility, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle compatibility: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	requestURL, err := url.JoinPath(
		c.baseURL,
		"/v1/accounts/vehicles",
		request.VIN,
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
	httpResponse, err := c.httpClient(cfg).Do(httpRequest)
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
	result := compatibilityResponseToProto(request.VIN, &responseBody)
	if rawStruct, err := compatibilityRawStruct(data); err == nil {
		result.SetRaw(rawStruct)
	}
	return result, nil
}
