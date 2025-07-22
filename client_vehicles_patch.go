package mbz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
)

// PatchVehiclesRequest is the request for [Client.PatchVehicles].
type PatchVehiclesRequest struct {
	// Vehicles is the list of vehicles to patch.
	Vehicles []vehiclesv1.Vehicle `json:"vehicles"`
}

// PatchVehiclesResponse is the response for [Client.PatchVehicles].
type PatchVehiclesResponse struct {
	// Vehicles is the list of vehicles returned by the API.
	Vehicles []vehiclesv1.Vehicle `json:"vehicles"`
}

// PatchVehicles patches vehicles. Only the deltaPush field is supported.
func (c *Client) PatchVehicles(ctx context.Context, request *PatchVehiclesRequest) (_ *PatchVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: patch vehicles: %w", err)
		}
	}()
	requestBodyData, err := json.Marshal(request.Vehicles)
	if err != nil {
		return nil, err
	}
	httpRequest, err := c.newRequest(ctx, http.MethodPatch, "/v1/accounts/vehicles", bytes.NewReader(requestBodyData))
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
	return &PatchVehiclesResponse{
		Vehicles: vehicles,
	}, nil
}
