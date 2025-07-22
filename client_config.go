package mbz

import (
	"log/slog"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ClientConfig configures a [Client].
type ClientConfig struct {
	region       Region
	tokenSource  oauth2.TokenSource
	oauth2Config *clientcredentials.Config
	logger       Logger
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

// WithOAuth2Config sets the OAuth2 config for the client.
func WithOAuth2Config(oauth2Config *clientcredentials.Config) ClientOption {
	return func(config *ClientConfig) {
		config.oauth2Config = oauth2Config
	}
}

// WithLogger sets the [Logger] for the [Client].
func WithLogger(logger Logger) ClientOption {
	return func(cc *ClientConfig) {
		cc.logger = logger
	}
}

// WithSlogLogger sets the [slog.Logger] for the [Client].
func WithSlogLogger(logger *slog.Logger) ClientOption {
	return func(cc *ClientConfig) {
		cc.logger = slogLogger{logger: logger}
	}
}
