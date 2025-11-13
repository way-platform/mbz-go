package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/adrg/xdg"
	"github.com/way-platform/mbz-go"
	"golang.org/x/oauth2"
)

// File storing authentication credentials for the CLI.
type File struct {
	// Region is the region of the credentials.
	Region mbz.Region `json:"region"`
	// Token is the current stored client credentials.
	Token *oauth2.Token `json:"token"`
	// APIKey is the API key for Vehicle Specification API.
	APIKey string `json:"apiKey,omitempty"`
}

func resolveFilepath() (string, error) {
	return xdg.ConfigFile("mbz-go/auth.json")
}

// NewOAuth2Client creates a new Mercedes-Benz API client using the current CLI credentials.
func NewOAuth2Client(ctx context.Context, opts ...mbz.ClientOption) (*mbz.Client, error) {
	cf, err := ReadFile()
	if err != nil {
		return nil, err
	}
	if cf.Token == nil || cf.Token.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("invalid token, please login using `mbz auth login`")
	}
	return mbz.NewClient(
		ctx,
		append(
			opts,
			mbz.WithRegion(cf.Region),
			mbz.WithOAuth2TokenSource(oauth2.StaticTokenSource(cf.Token)),
		)...,
	)
}

// NewClientWithAPIKey creates a new Mercedes-Benz API client using the API key from the CLI credentials.
func NewClientWithAPIKey(ctx context.Context, opts ...mbz.ClientOption) (*mbz.Client, error) {
	cf, err := ReadFile()
	if err != nil {
		return nil, err
	}
	if cf.APIKey == "" {
		return nil, fmt.Errorf("no API key found, please login using `mbz auth login --api-key <api-key>`")
	}
	return mbz.NewClient(
		ctx,
		append(
			opts,
			mbz.WithAPIKey(cf.APIKey),
		)...,
	)
}

// ReadFile reads the currently stored [File].
func ReadFile() (*File, error) {
	fp, err := resolveFilepath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fp); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no credentials found, please login using `mbz auth login`")
		}
		return nil, err
	}
	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// writeFile writes the stored [File].
func writeFile(f *File) error {
	fp, err := resolveFilepath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fp, data, 0o600)
}

// removeFile removes the stored [File].
func removeFile() error {
	fp, err := resolveFilepath()
	if err != nil {
		return err
	}
	return os.RemoveAll(fp)
}
