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

// PatchVehicles patches vehicles. Only the delta_push field is supported.
func (c *Client) PatchVehicles(
	ctx context.Context,
	request *fleetv1.PatchVehiclesRequest,
) (_ *fleetv1.PatchVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: patch vehicles: %w", err)
		}
	}()
	apiVehicles := make([]vehiclesv1.Vehicle, 0, len(request.GetVehicles()))
	for _, v := range request.GetVehicles() {
		av := vehiclesv1.Vehicle{VIN: v.GetVin()}
		if v.HasDeltaPush() {
			dp := v.GetDeltaPush()
			av.DeltaPush = &dp
		}
		apiVehicles = append(apiVehicles, av)
	}
	requestBodyData, err := json.Marshal(apiVehicles)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	requestURL, err := url.JoinPath(c.baseURL, "/v1/accounts/vehicles")
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		requestURL,
		bytes.NewReader(requestBodyData),
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accept", "application/json")
	httpResponse, err := c.httpClient(c.config).Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if closeErr := httpResponse.Body.Close(); closeErr != nil {
			log.Printf("mbz: failed to close response body: %v", closeErr)
		}
	}()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	return &fleetv1.PatchVehiclesResponse{}, nil
}
