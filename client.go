package mbz

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"golang.org/x/oauth2"
)

// Client to the Mercedes-Benz management APIs.
type Client struct {
	baseURL string
	config  clientConfig
}

// NewClient creates a new Mercedes-Benz API client.
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	config := newClientConfig()
	for _, opt := range opts {
		opt(&config)
	}
	client := &Client{
		config: config,
	}
	switch config.region {
	case RegionECE:
		client.baseURL = BaseURLECE
	case RegionAMAPNA:
		client.baseURL = BaseURLAMAPNA
	default:
		return nil, fmt.Errorf("invalid region: %s", config.region)
	}
	// Initialize OAuth2 token source if credentials are provided
	if config.clientID != "" && config.clientSecret != "" {
		oauth2Config, err := NewOAuth2Config(config.region, config.clientID, config.clientSecret)
		if err != nil {
			return nil, err
		}
		client.config.tokenSource = oauth2Config.TokenSource(ctx)
	}
	return client, nil
}

// clientConfig configures a [Client].
type clientConfig struct {
	region       Region
	clientID     string
	clientSecret string
	apiKey       string
	tokenSource  oauth2.TokenSource
	debug        bool
	retryCount   int
	timeout      time.Duration
	interceptors []func(http.RoundTripper) http.RoundTripper
}

func newClientConfig() clientConfig {
	return clientConfig{
		region:     RegionECE,
		retryCount: 3,
		timeout:    30 * time.Second,
	}
}

func (cc clientConfig) with(opts ...ClientOption) clientConfig {
	for _, opt := range opts {
		opt(&cc)
	}
	return cc
}

// ClientOption is a configuration option for a [Client].
type ClientOption func(*clientConfig)

// WithRegion sets the region for the client.
func WithRegion(region Region) ClientOption {
	return func(config *clientConfig) {
		config.region = region
	}
}

// WithOAuth2TokenSource sets the OAuth2 token source for the client.
func WithOAuth2TokenSource(tokenSource oauth2.TokenSource) ClientOption {
	return func(config *clientConfig) {
		config.tokenSource = tokenSource
	}
}

// WithClientID sets the client ID for OAuth2 authentication.
func WithClientID(clientID string) ClientOption {
	return func(config *clientConfig) {
		config.clientID = clientID
	}
}

// WithClientSecret sets the client secret for OAuth2 authentication.
func WithClientSecret(clientSecret string) ClientOption {
	return func(config *clientConfig) {
		config.clientSecret = clientSecret
	}
}

// WithAPIKey sets the API key for the client.
func WithAPIKey(apiKey string) ClientOption {
	return func(config *clientConfig) {
		config.apiKey = apiKey
	}
}

// WithDebug toggles debug mode (request/response dumps to stderr).
func WithDebug(debug bool) ClientOption {
	return func(config *clientConfig) {
		config.debug = debug
	}
}

// WithRetryCount sets the number of retries for API requests.
func WithRetryCount(retryCount int) ClientOption {
	return func(config *clientConfig) {
		config.retryCount = retryCount
	}
}

// WithTimeout sets the timeout for API requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(config *clientConfig) {
		config.timeout = timeout
	}
}

// WithInterceptor adds a request interceptor for the [Client].
func WithInterceptor(interceptor func(http.RoundTripper) http.RoundTripper) ClientOption {
	return func(config *clientConfig) {
		config.interceptors = append(config.interceptors, interceptor)
	}
}

func (c *Client) httpClient(cfg clientConfig) *http.Client {
	transport := http.DefaultTransport
	if cfg.debug {
		transport = &debugTransport{next: transport}
	}
	if cfg.tokenSource != nil {
		transport = &oauth2.Transport{
			Source: cfg.tokenSource,
			Base:   transport,
		}
	}
	if cfg.apiKey != "" {
		transport = &apiKeyTransport{
			apiKey: cfg.apiKey,
			next:   transport,
		}
	}
	if len(cfg.interceptors) > 0 {
		transport = &interceptorTransport{
			interceptors: cfg.interceptors,
			next:         transport,
		}
	}
	if cfg.retryCount > 0 {
		transport = &retryTransport{
			maxRetries: cfg.retryCount,
			next:       transport,
		}
	}
	return &http.Client{
		Timeout:   cfg.timeout,
		Transport: transport,
	}
}

func getUserAgent() string {
	userAgent := "WayPlatformMBZGo"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		userAgent += "/" + info.Main.Version
	}
	return userAgent
}
