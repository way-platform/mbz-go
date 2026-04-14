package mbz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
	fleetv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/fleet/v1"
)

// PostVehicleServices activates or deactivates data services for vehicles.
func (c *Client) PostVehicleServices(
	ctx context.Context,
	request *fleetv1.PostVehicleServicesRequest,
) (_ *fleetv1.PostVehicleServicesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: post vehicle services: %w", err)
		}
	}()
	var requestBody []vehiclesv1.DesiredServiceStatusRequest
	for _, input := range request.GetVehicleServiceInputs() {
		services := make([]vehiclesv1.DesiredServiceStatus, 0, len(input.GetServices()))
		for _, svc := range input.GetServices() {
			services = append(services, vehiclesv1.DesiredServiceStatus{
				ServiceID:     svc.GetServiceId(),
				DesiredStatus: vehiclesv1.Status(svc.GetDesiredStatus()),
			})
		}
		requestBody = append(requestBody, vehiclesv1.DesiredServiceStatusRequest{
			VIN:      input.GetVin(),
			Services: services,
		})
	}
	requestBodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	requestURL, err := url.JoinPath(c.baseURL, "/v2/accounts/vehicles/services")
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		requestURL,
		bytes.NewReader(requestBodyData),
	)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.Header.Set("Content-Type", "application/json")
	httpResponse, err := c.httpClient(c.config).Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := httpResponse.Body.Close(); closeErr != nil {
			log.Printf("mbz: failed to close response body: %v", closeErr)
		}
	}()
	if httpResponse.StatusCode != http.StatusAccepted {
		return nil, newResponseError(httpResponse)
	}
	return &fleetv1.PostVehicleServicesResponse{}, nil
}
