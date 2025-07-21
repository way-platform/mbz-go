package mbz

import (
	"context"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Client to the Mercedes-Benz management APIs.
type Client struct {
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

// ClientConfig configures a [Client].
type ClientConfig struct {
	region       Region
	tokenSource  oauth2.TokenSource
	oauth2Config *clientcredentials.Config
}

// ClientOption is a configuration option for a [Client].
type ClientOption func(*ClientConfig)

// WithRegion sets the region for the client.
func WithRegion(region Region) ClientOption {
	return func(config *ClientConfig) {
		config.region = region
	}
}

// WithOAuth2TokenSource sets the OAuth2 token source for the client.
func WithOAuth2TokenSource(tokenSource oauth2.TokenSource) ClientOption {
	return func(config *ClientConfig) {
		config.tokenSource = tokenSource
	}
}

func WithOAuth2Config(oauth2Config *clientcredentials.Config) ClientOption {
	return func(config *ClientConfig) {
		config.oauth2Config = oauth2Config
	}
}
