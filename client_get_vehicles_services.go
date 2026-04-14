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

// GetVehicleServices gets the actual service status for a vehicle.
func (c *Client) GetVehicleServices(
	ctx context.Context,
	request *fleetv1.GetVehicleServicesRequest,
) (_ *fleetv1.GetVehicleServicesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle services: %w", err)
		}
	}()
	requestURL, err := url.JoinPath(c.baseURL, "/v2/accounts/vehicles", request.GetVin(), "services")
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
	var responseBody vehiclesv1.VehicleServiceStatus
	if err := json.Unmarshal(data, &responseBody); err != nil {
		return nil, err
	}
	services := make([]*fleetv1.ServiceStatus, 0, len(responseBody.Services))
	for _, s := range responseBody.Services {
		ps := &fleetv1.ServiceStatus{}
		ps.SetServiceId(s.ServiceID)
		ps.SetStatus(string(s.Status))
		ps.SetDesiredStatus(string(s.DesiredStatus))
		ps.SetOrderTime(s.OrderTime)
		services = append(services, ps)
	}
	resp := &fleetv1.GetVehicleServicesResponse{}
	resp.SetDeltaPush(responseBody.DeltaPush)
	resp.SetServices(services)
	return resp, nil
}
