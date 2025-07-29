package mbz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
)

// PatchVehiclesRequest is the request for [Client.PatchVehicles].
type PatchVehiclesRequest struct {
	// Vehicles is the list of vehicles to patch.
	Vehicles []vehiclesv1.Vehicle `json:"vehicles"`
}

// PatchVehiclesResponse is the response for [Client.PatchVehicles].
type PatchVehiclesResponse struct{}

// PatchVehicles patches vehicles. Only the deltaPush field is supported.
func (c *Client) PatchVehicles(ctx context.Context, request *PatchVehiclesRequest) (_ *PatchVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: patch vehicles: %w", err)
		}
	}()
	requestBodyData, err := json.Marshal(request.Vehicles)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	httpRequest, err := c.newRequest(ctx, http.MethodPatch, "/v1/accounts/vehicles", bytes.NewReader(requestBodyData))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accept", "application/json")
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	return &PatchVehiclesResponse{}, nil
}
