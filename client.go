package mbz

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"runtime/debug"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/oauth2"
)

// Client to the Mercedes-Benz management APIs.
type Client struct {
	baseURL    string
	httpClient *retryablehttp.Client
	config     ClientConfig
}

// NewClient creates a new Mercedes-Benz API client.
func NewClient(opts ...ClientOption) *Client {
	config := ClientConfig{}
	for _, opt := range opts {
		opt(&config)
	}
	client := &Client{
		httpClient: retryablehttp.NewClient(),
		config:     config,
	}
	client.httpClient.Logger = noopLogger{}
	switch config.region {
	case RegionECE:
		client.baseURL = BaseURLECE
	case RegionAMAPNA:
		client.baseURL = BaseURLAMAPNA
	default: // default to ECE region
		client.baseURL = BaseURLECE
	}
	switch {
	case config.oauth2Config != nil:
		client.httpClient.HTTPClient = config.oauth2Config.Client(context.Background())
	case config.tokenSource != nil:
		client.httpClient.HTTPClient.Transport = &oauth2.Transport{
			Source: config.tokenSource,
			Base:   client.httpClient.HTTPClient.Transport,
		}
	}
	return client
}

func (c *Client) newRequest(ctx context.Context, method, requestPath string, body io.Reader) (_ *retryablehttp.Request, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("new request: %w", err)
		}
	}()
	requestURL, err := url.JoinPath(c.baseURL, requestPath)
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}
	request, err := retryablehttp.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", getUserAgent())
	return request, nil
}

func getUserAgent() string {
	userAgent := "WayPlatformMBZGo"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		userAgent += "/" + info.Main.Version
	}
	return userAgent
}
