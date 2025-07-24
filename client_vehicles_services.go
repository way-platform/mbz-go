package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	// Services with the service availability.
	Services []vehiclesv1.ServiceStatus `json:"services"`
}

// GetVehicleServices gets the actual service status for a vehicle.
func (c *Client) GetVehicleServices(ctx context.Context, request *GetVehicleServicesRequest) (_ *GetVehicleServicesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle services: %w", err)
		}
	}()
	httpRequest, err := c.newRequest(ctx, http.MethodGet, "/v1/accounts/vehicles/"+request.VIN+"/services", nil)
	if err != nil {
		return nil, err
	}
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
	log.Println(string(data))
	var responseBody []vehiclesv1.ServiceStatus
	if err := json.Unmarshal(data, &responseBody); err != nil {
		return nil, err
	}
	return &GetVehicleServicesResponse{
		Services: responseBody,
	}, nil
}
