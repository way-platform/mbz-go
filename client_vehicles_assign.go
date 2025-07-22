package mbz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
)

// AssignVehiclesRequest is the request for [Client.AssignVehicles].
type AssignVehiclesRequest struct {
	// VINs is the list of VINs to assign to your account.
	VINs []string `json:"vins"`
}

// AssignVehiclesResponse is the response for [Client.AssignVehicles].
type AssignVehiclesResponse struct{}

// AssignVehicles lists the vehicles for the current account.
func (c *Client) AssignVehicles(ctx context.Context, request *AssignVehiclesRequest) (_ *AssignVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: assign vehicles: %w", err)
		}
	}()
	requestBody := make([]vehiclesv1.Vehicle, 0, len(request.VINs))
	for _, vin := range request.VINs {
		requestBody = append(requestBody, vehiclesv1.Vehicle{
			VIN: vin,
		})
	}
	requestBodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	httpRequest, err := c.newRequest(ctx, http.MethodPost, "/v1/accounts/vehicles", bytes.NewReader(requestBodyData))
	if err != nil {
		return nil, err
	}
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusCreated {
		return nil, newResponseError(httpResponse)
	}
	return &AssignVehiclesResponse{}, nil
}
