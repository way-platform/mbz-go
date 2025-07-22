package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
func (c *Client) ListVehicles(ctx context.Context, request *ListVehiclesRequest) (_ *ListVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: list vehicles: %w", err)
		}
	}()
	httpRequest, err := c.newRequest(ctx, http.MethodGet, "/v1/accounts/vehicles", nil)
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
	var vehicles []vehiclesv1.Vehicle
	if err := json.Unmarshal(data, &vehicles); err != nil {
		return nil, err
	}
	return &ListVehiclesResponse{
		Vehicles: vehicles,
	}, nil
}
