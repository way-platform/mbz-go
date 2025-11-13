package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
)

// GetVehicleServicesRequest is the request for [Client.GetVehicleServices].
type GetVehicleServicesRequest struct {
	// VIN of the vehicle to get the services for.
	VIN string `json:"vin"`
}

// GetVehicleServicesResponse is the response for [Client.GetVehicleServices].
type GetVehicleServicesResponse struct {
	// DeltaPush indicates if delta push is enabled for the vehicle.
	DeltaPush bool `json:"deltaPush"`
	// Services with the service availability.
	Services []vehiclesv1.ServiceStatus `json:"services"`
}

// GetVehicleServices gets the actual service status for a vehicle.
func (c *Client) GetVehicleServices(ctx context.Context, request *GetVehicleServicesRequest, opts ...ClientOption) (_ *GetVehicleServicesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle services: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	httpRequest, err := c.newRequest(ctx, http.MethodGet, "/v2/accounts/vehicles/"+request.VIN+"/services", nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := c.httpClient(cfg).Do(httpRequest)
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
	var responseBody vehiclesv1.VehicleServiceStatus
	if err := json.Unmarshal(data, &responseBody); err != nil {
		return nil, err
	}
	return &GetVehicleServicesResponse{
		DeltaPush: responseBody.DeltaPush,
		Services:  responseBody.Services,
	}, nil
}
