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

// ListVehicles lists the vehicles for the current account.
func (c *Client) ListVehicles(
	ctx context.Context,
	request *fleetv1.ListVehiclesRequest,
) (_ *fleetv1.ListVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: list vehicles: %w", err)
		}
	}()
	requestURL, err := url.JoinPath(c.baseURL, "/v1/accounts/vehicles")
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
	var apiVehicles []vehiclesv1.Vehicle
	if err := json.Unmarshal(data, &apiVehicles); err != nil {
		return nil, err
	}
	vehicles := make([]*fleetv1.Vehicle, 0, len(apiVehicles))
	for _, v := range apiVehicles {
		pv := &fleetv1.Vehicle{}
		pv.SetVin(v.VIN)
		if v.DeltaPush != nil {
			pv.SetDeltaPush(*v.DeltaPush)
		}
		vehicles = append(vehicles, pv)
	}
	resp := &fleetv1.ListVehiclesResponse{}
	resp.SetVehicles(vehicles)
	return resp, nil
}
