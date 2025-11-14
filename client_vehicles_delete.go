package mbz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclesv1"
)

// DeleteVehiclesRequest is the request for [Client.DeleteVehicles].
type DeleteVehiclesRequest struct {
	// VINs is the list of VINs to delete from your account.
	VINs []string `json:"vins"`
}

// DeleteVehiclesResponse is the response for [Client.DeleteVehicles].
type DeleteVehiclesResponse struct{}

// DeleteVehicles lists the vehicles for the current account.
func (c *Client) DeleteVehicles(
	ctx context.Context,
	request *DeleteVehiclesRequest,
	opts ...ClientOption,
) (_ *DeleteVehiclesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: delete vehicles: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
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
	requestURL, err := url.JoinPath(c.baseURL, "/v1/accounts/vehicles")
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL, bytes.NewReader(requestBodyData))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.Header.Set("Content-Type", "application/json")
	httpResponse, err := c.httpClient(cfg).Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	return &DeleteVehiclesResponse{}, nil
}
