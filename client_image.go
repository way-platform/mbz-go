package mbz

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
)

// GetImageRequest is the request for [Client.GetImage].
type GetImageRequest struct {
	// ImageID is the UUID of the image to download.
	ImageID string `json:"imageId"`
}

// GetImageResponse is the response for [Client.GetImage].
type GetImageResponse struct {
	// Data is the binary image data.
	Data []byte
	// ContentType is the MIME type of the image (e.g., "image/png", "image/jpeg", "image/webp").
	ContentType string
}

// GetImage downloads the vehicle image for a given image ID.
// This method requires API key authentication via [WithAPIKey].
func (c *Client) GetImage(
	ctx context.Context,
	request *GetImageRequest,
	opts ...ClientOption,
) (_ *GetImageResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get image: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	if request.ImageID == "" {
		return nil, fmt.Errorf("image ID is required")
	}
	requestURL := fmt.Sprintf("%s/images/%s", vehiclespecificationfleetv1.BaseURL, request.ImageID)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	// Accept any image format
	httpRequest.Header.Set("Accept", "image/png,image/jpeg,image/webp,*/*")
	httpResponse, err := c.httpClient(cfg).Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
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
	return &GetImageResponse{
		Data:        imageData,
		ContentType: contentType,
	}, nil
}

