package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/way-platform/mbz-go/api/servicesv1"
)

// ListServicesRequest is the request for [Client.ListServices].
type ListServicesRequest struct {
	// Details is a flag to include service details in the response.
	Details bool `json:"details"`
}

// ListServicesResponse is the response for [Client.ListServices].
type ListServicesResponse struct {
	// Services is the list of services returned by the API.
	Services []servicesv1.Service `json:"services"`
	// Version is the version of the services spec.
	Version string `json:"version,omitempty"`
}

// ListServices lists the vehicles for the current account.
func (c *Client) ListServices(ctx context.Context, request *ListServicesRequest) (_ *ListServicesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: list vehicles: %w", err)
		}
	}()
	url := "/v2/accounts/services"
	if request.Details {
		url += "/v2/accounts/services/details"
	}
	httpRequest, err := c.newRequest(ctx, http.MethodGet, url, nil)
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
	var responseBody servicesv1.GetAllServicesWithSignalsAndCommandsResponse
	if err := json.Unmarshal(data, &responseBody); err != nil {
		return nil, err
	}
	return &ListServicesResponse{
		Services: responseBody.Services,
		Version:  responseBody.Version,
	}, nil
}
