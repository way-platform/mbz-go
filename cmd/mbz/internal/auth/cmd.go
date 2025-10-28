package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/way-platform/mbz-go"
	"golang.org/x/oauth2"
	"golang.org/x/term"
)

// NewOAuth2Client creates a new Mercedes-Benz API client using the current CLI credentials.
func NewOAuth2Client() (*mbz.Client, error) {
	cf, err := ReadFile()
	if err != nil {
		return nil, err
	}
	return mbz.NewClient(
		mbz.WithRegion(cf.Region),
		mbz.WithOAuth2TokenSource(oauth2.StaticTokenSource(&cf.Credentials)),
		mbz.WithSlogLogger(slog.Default()),
	)
}

// NewClientWithAPIKey creates a new Mercedes-Benz API client using the API key from the CLI credentials.
func NewClientWithAPIKey() (*mbz.Client, error) {
	cf, err := ReadFile()
	if err != nil {
		return nil, err
	}
	if cf.APIKey == "" {
		return nil, fmt.Errorf("no API key found, please login using `mbz auth login apikey`")
	}
	return mbz.NewClient(
		mbz.WithAPIKey(cf.APIKey),
		mbz.WithSlogLogger(slog.Default()),
	)
}

// NewCommand returns a new [cobra.Command] for CLI authentication.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate to the Mercedes-Benz API",
	}
	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newLogoutCommand())
	return cmd
}

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the Mercedes-Benz API",
		Long:  "Login to the Mercedes-Benz API. Defaults to OAuth login if no subcommand is specified.",
	}
	oauthCmd := newOAuthLoginCommand()
	cmd.AddCommand(oauthCmd)
	cmd.AddCommand(newAPIKeyLoginCommand())
	cmd.RunE = oauthCmd.RunE
	cmd.Flags().AddFlagSet(oauthCmd.Flags())
	return cmd
}

func newOAuthLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oauth",
		Short: "Login with OAuth2 credentials",
	}
	region := cmd.Flags().String("region", string(mbz.RegionECE), "region to use for authentication")
	clientID := cmd.Flags().String("client-id", "-", "client ID to use for authentication")
	clientSecret := cmd.Flags().String("client-secret", "-", "client secret to use for authentication")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if *clientID == "-" {
			cmd.Println("\nEnter OAuth2 client ID:")
			input, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			*clientID = string(input)
		}
		if *clientSecret == "-" {
			cmd.Println("\nEnter OAuth2 client secret:")
			input, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			*clientSecret = string(input)
		}
		config, err := mbz.NewOAuth2Config(mbz.Region(*region), *clientID, *clientSecret)
		if err != nil {
			return err
		}
		token, err := config.Token(cmd.Context())
		if err != nil {
			return err
		}
		if err := writeFile(&File{
			Region:      mbz.Region(*region),
			Credentials: *token,
		}); err != nil {
			return err
		}
		cmd.Printf("Logged in to %s.\n", *region)
		return nil
	}
	return cmd
}

func newAPIKeyLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apikey",
		Short: "Login with API key",
	}
	apiKey := cmd.Flags().String("api-key", "-", "API key to use for authentication")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if *apiKey == "-" {
			cmd.Println("\nEnter API key:")
			input, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			*apiKey = string(input)
		}
		if err := writeFile(&File{
			APIKey: *apiKey,
		}); err != nil {
			return err
		}
		cmd.Println("Logged in with API key")
		return nil
	}
	return cmd
}

func newLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from the Mercedes-Benz API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := removeFile(); err != nil {
				return err
			}
			cmd.Println("Logged out.")
			return nil
		},
	}
}

// File storing authentication credentials for the CLI.
type File struct {
	// Region is the region of the credentials.
	Region mbz.Region `json:"region"`
	// Credentials is the current stored client credentials.
	Credentials oauth2.Token `json:"clientCredentials"`
	// APIKey is the API key for Vehicle Specification API.
	APIKey string `json:"apiKey,omitempty"`
}

func (cf *File) isExpired() bool {
	return cf.Credentials.Expiry.Before(time.Now()) && cf.APIKey == ""
}

func resolveFilepath() (string, error) {
	return xdg.ConfigFile("mbz-go/auth.json")
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
	if f.isExpired() {
		return nil, fmt.Errorf("credentials expired, please login again")
	}
	return &f, nil
}

// writeFile writes the stored [credentialsFile].
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
