package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
)

// GetVehicleSpecificationRequest is the request for [Client.GetVehicleSpecification].
type GetVehicleSpecificationRequest struct {
	// VIN is the VIN (or FIN) of the vehicle (17 characters).
	VIN string `json:"vin"`
	// Locale is the market locale.
	Locale string `json:"locale"`
}

// GetVehicleSpecification gets the vehicle marketing information for a given vehicle ID.
// This method requires API key authentication via [WithAPIKey].
func (c *Client) GetVehicleSpecification(
	ctx context.Context,
	request *GetVehicleSpecificationRequest,
	opts ...ClientOption,
) (_ *mbzv1.VehicleSpecification, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle specification: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	if request.VIN == "" {
		return nil, fmt.Errorf("VIN is required")
	}
	values := url.Values{}
	if request.Locale != "" {
		values.Set("locale", request.Locale)
	} else {
		values.Set("locale", string(vehiclespecificationfleetv1.LocalesEnUS))
	}
	// Set all optional parameters to true to maximize data retrieval
	values.Set("additionalSpecs", "true")
	values.Set("optionsNullDescription", "true")
	values.Set("options", "true")
	values.Set("technicalData", "true")
	values.Set("payloadNullValues", "true")
	requestURL := fmt.Sprintf("%s/vehicles/%s", vehiclespecificationfleetv1.BaseURL, request.VIN)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.URL.RawQuery = values.Encode()
	httpRequest.Header.Set("Accept", "application/json")
	httpResponse, err := c.httpClient(cfg).Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	responseData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var openAPIResp vehiclespecificationfleetv1.VehicleSpecificationResponse
	if err := json.Unmarshal(responseData, &openAPIResp); err != nil {
		return nil, err
	}
	return vehicleDataToProto(openAPIResp.VehicleData), nil
}
