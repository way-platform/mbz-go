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

// AssignVehicles assigns vehicles to the current account.
func (c *Client) AssignVehicles(
	ctx context.Context,
	request *fleetv1.AssignVehiclesRequest,
) (_ *fleetv1.AssignVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: assign vehicles: %w", err)
		}
	}()
	requestBody := make([]vehiclesv1.Vehicle, 0, len(request.GetVins()))
	for _, vin := range request.GetVins() {
		requestBody = append(requestBody, vehiclesv1.Vehicle{
			VIN: vin,
		})
	}
	requestBodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	requestURL, err := url.JoinPath(c.baseURL, "/v1/accounts/vehicles")
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
	if httpResponse.StatusCode != http.StatusCreated {
		return nil, newResponseError(httpResponse)
	}
	return &fleetv1.AssignVehiclesResponse{}, nil
}
