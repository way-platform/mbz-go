package mbz

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	vehiclespecv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/vehiclespec/v1"
)

// GetImage downloads the vehicle image for a given image ID.
// This method requires API key authentication via [WithAPIKey].
func (c *Client) GetImage(
	ctx context.Context,
	request *vehiclespecv1.GetImageRequest,
) (_ *vehiclespecv1.GetImageResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get image: %w", err)
		}
	}()
	if request.GetImageId() == "" {
		return nil, fmt.Errorf("image ID is required")
	}
	requestURL := fmt.Sprintf("%s/images/%s", vehiclespecificationfleetv1.BaseURL, request.GetImageId())
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	// Accept any image format.
	httpRequest.Header.Set("Accept", "image/png,image/jpeg,image/webp,*/*")
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
	imageData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	contentType := httpResponse.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	resp := &vehiclespecv1.GetImageResponse{}
	resp.SetData(imageData)
	resp.SetContentType(contentType)
	return resp, nil
}
