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

// GetVehicleImageIdsRequest is the request for [Client.GetVehicleImageIds].
type GetVehicleImageIdsRequest struct {
	// VIN is the VIN (or FIN) of the vehicle (17 characters).
	VIN string `json:"vin"`
	// Background is an optional property that defines the images background.
	// The default value is false (transparent). Set to true to retrieve images
	// with a high level of detail and realistic reflections and light incidence.
	Background bool `json:"background,omitempty"`
	// FileFormat is an optional property that defines the images format.
	// Valid values: png, jpeg, webp. Default is webp.
	FileFormat string `json:"fileFormat,omitempty"`
}

// GetVehicleImageIds gets the vehicle image IDs for a given vehicle ID.
// This method requires API key authentication via [WithAPIKey].
func (c *Client) GetVehicleImageIds(
	ctx context.Context,
	request *GetVehicleImageIdsRequest,
	opts ...ClientOption,
) (_ *mbzv1.VehicleImages, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle image ids: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	if request.VIN == "" {
		return nil, fmt.Errorf("VIN is required")
	}
	values := url.Values{}
	if request.Background {
		values.Set("background", "true")
	}
	if request.FileFormat != "" {
		values.Set("fileFormat", request.FileFormat)
	}
	requestURL := fmt.Sprintf("%s/vehicle-images/%s", vehiclespecificationfleetv1.BaseURL, request.VIN)
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
	var openAPIResp vehiclespecificationfleetv1.VehicleImagesResponse
	if err := json.Unmarshal(responseData, &openAPIResp); err != nil {
		return nil, err
	}
	return vehicleImagesResponseToProto(&openAPIResp), nil
}

func vehicleImagesResponseToProto(
	response *vehiclespecificationfleetv1.VehicleImagesResponse,
) *mbzv1.VehicleImages {
	var output mbzv1.VehicleImages
	if response == nil {
		return &output
	}
	// Set image ID fields
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
