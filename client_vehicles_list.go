package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
)

// ListVehiclesRequest is the request for [Client.ListVehicles].
type ListVehiclesRequest struct{}

// ListVehiclesResponse is the response for [Client.ListVehicles].
type ListVehiclesResponse struct {
	// Vehicles is the list of vehicles returned by the API.
	Vehicles []vehiclesv1.Vehicle `json:"vehicles"`
}

// ListVehicles lists the vehicles for the current account.
func (c *Client) ListVehicles(ctx context.Context, request *ListVehiclesRequest, opts ...ClientOption) (_ *ListVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: list vehicles: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	requestURL, err := url.JoinPath(c.baseURL, "/v1/accounts/vehicles")
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var vehicles []vehiclesv1.Vehicle
	if err := json.Unmarshal(data, &vehicles); err != nil {
		return nil, err
	}
	return &ListVehiclesResponse{
		Vehicles: vehicles,
	}, nil
}
