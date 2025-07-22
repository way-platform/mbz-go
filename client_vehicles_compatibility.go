package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
)

// GetVehicleCompatibilityRequest is the request for [Client.GetVehicleCompatibility].
type GetVehicleCompatibilityRequest struct {
	// VIN of the vehicle to get the compatibility for.
	VIN string
}

// GetVehicleCompatibilityResponse is the response for [Client.GetVehicleCompatibility].
type GetVehicleCompatibilityResponse struct {
	// VIN of the requested vehicle.
	VIN string `json:"vin"`

	// VehicleType is the type of the requested vehicle.
	VehicleType string `json:"vehicleType,omitempty"`

	// VehicleProvidesConnectivity indicates the base compatibility to data-services for the requested vehicle.
	VehicleProvidesConnectivity bool `json:"vehicleProvidesConnectivity"`

	// Services with the service availability.
	Services []vehiclesv1.CompatibilityGenericService `json:"services"`
}

// GetVehicleCompatibility gets the compatibility of a vehicle.
func (c *Client) GetVehicleCompatibility(ctx context.Context, request *GetVehicleCompatibilityRequest) (_ *GetVehicleCompatibilityResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle compatibility: %w", err)
		}
	}()
	httpRequest, err := c.newRequest(ctx, http.MethodGet, "/v1/accounts/vehicles/"+request.VIN+"/compatibilites", nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
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
	return &GetVehicleCompatibilityResponse{
		VIN:                         responseBody.VIN,
		VehicleType:                 responseBody.VehicleType,
		VehicleProvidesConnectivity: responseBody.VehicleProvidesConnectivity,
		Services:                    responseBody.Services,
	}, nil
}
