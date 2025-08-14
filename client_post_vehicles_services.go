package mbz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
)

type DesiredStatus string

const (
	DesiredStatusActive   DesiredStatus = "ACTIVE"
	DesiredStatusInactive DesiredStatus = "INACTIVE"
)

type VehicleServices struct {
	ServiceID     string        `json:"serviceId"`
	DesiredStatus DesiredStatus `json:"desiredStatus"`
}

// PostVehicleServicesRequest is the request for [Client.PostVehicleServices].
type PostVehicleServicesRequest struct {
	DesiredServiceStatusInput []DesiredServiceStatusInput `json:"desiredServiceStatusInput"`
}

type DesiredServiceStatusInput struct {
	// VIN of the vehicle to get the services for.
	VIN string `json:"vin"`
	// Services is the list of services to activate or deactivate for the given VIN.
	Services []VehicleServices `json:"services"`
}

// PostVehicleServicesResponse is the response for [Client.PostVehicleServices].
type PostVehicleServicesResponse struct{}

// PostVehicleServices posts the actual service status for a vehicle.
func (c *Client) PostVehicleServices(
	ctx context.Context,
	request *PostVehicleServicesRequest,
) (_ *PostVehicleServicesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: post vehicle services: %w", err)
		}
	}()
	var requestBody []vehiclesv1.DesiredServiceStatusRequest
	for _, desiredServiceStatusInput := range request.DesiredServiceStatusInput {
		services := make([]vehiclesv1.DesiredServiceStatus, 0, len(desiredServiceStatusInput.Services))
		for _, service := range desiredServiceStatusInput.Services {
			services = append(services, vehiclesv1.DesiredServiceStatus{
				ServiceID:     service.ServiceID,
				DesiredStatus: vehiclesv1.Status(service.DesiredStatus),
			})
		}
		requestBody = append(requestBody, vehiclesv1.DesiredServiceStatusRequest{
			VIN:      desiredServiceStatusInput.VIN,
			Services: services,
		})
	}
	requestBodyData, err := json.Marshal(requestBody)
	fmt.Println("requestBody", string(requestBodyData))
	if err != nil {
		return nil, err
	}
	httpRequest, err := c.newRequest(
		ctx,
		http.MethodPost,
		"/v2/accounts/vehicles/services",
		bytes.NewReader(requestBodyData),
	)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusAccepted {
		return nil, newResponseError(httpResponse)
	}
	return &PostVehicleServicesResponse{}, nil
}
