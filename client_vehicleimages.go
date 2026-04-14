package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
	vehiclespecv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/vehiclespec/v1"
)

// GetVehicleImageIds gets the vehicle image IDs for a given vehicle.
// This method requires API key authentication via [WithAPIKey].
func (c *Client) GetVehicleImageIds(
	ctx context.Context,
	request *vehiclespecv1.GetVehicleImageIdsRequest,
) (_ *vehiclespecv1.GetVehicleImageIdsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle image ids: %w", err)
		}
	}()
	if request.GetVin() == "" {
		return nil, fmt.Errorf("VIN is required")
	}
	values := url.Values{}
	if request.GetBackground() {
		values.Set("background", "true")
	}
	if request.GetFileFormat() != "" {
		values.Set("fileFormat", request.GetFileFormat())
	}
	requestURL := fmt.Sprintf(
		"%s/vehicle-images/%s",
		vehiclespecificationfleetv1.BaseURL,
		request.GetVin(),
	)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.URL.RawQuery = values.Encode()
	httpRequest.Header.Set("Accept", "application/json")
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
	responseData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var openAPIResp vehiclespecificationfleetv1.VehicleImagesResponse
	if err := json.Unmarshal(responseData, &openAPIResp); err != nil {
		return nil, err
	}
	resp := &vehiclespecv1.GetVehicleImageIdsResponse{}
	resp.SetVehicleImages(vehicleImagesResponseToProto(&openAPIResp))
	return resp, nil
}

func vehicleImagesResponseToProto(
	response *vehiclespecificationfleetv1.VehicleImagesResponse,
) *mbzv1.VehicleImages {
	var output mbzv1.VehicleImages
	if response == nil {
		return &output
	}
	if response.EXT150 != "" {
		output.SetExt150(response.EXT150)
	}
	if response.EXT330 != "" {
		output.SetExt330(response.EXT330)
	}
	if response.INT1 != "" {
		output.SetInt1(response.INT1)
	}
	return &output
}
