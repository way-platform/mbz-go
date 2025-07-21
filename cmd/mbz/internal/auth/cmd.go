package auth

import (
	"encoding/json"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/way-platform/mbz-go"
	"golang.org/x/oauth2"
)

// NewClient creates a new Mercedes-Benz Management API client using the current CLI credentials.
func NewClient() (*mbz.Client, error) {
	cf, err := readCredentialsFile()
	if err != nil {
		return nil, err
	}
	return mbz.NewClient(
		mbz.WithRegion(cf.Region),
		mbz.WithOAuth2TokenSource(oauth2.StaticTokenSource(&cf.Credentials)),
	), nil
}

// NewCommand returns a new [cobra.Command] for CLI authentication.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with the Mercedes-Benz Management API",
	}
	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newLogoutCommand())
	return cmd
}

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the Mercedes-Benz Management API",
	}
	region := cmd.Flags().String("region", string(mbz.RegionECE), "region to use for authentication")
	clientID := cmd.Flags().String("client-id", "", "client ID to use for authentication")
	cmd.MarkFlagRequired("client-id")
	clientSecret := cmd.Flags().String("client-secret", "", "client secret to use for authentication")
	cmd.MarkFlagRequired("client-secret")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		config, err := mbz.NewOAuth2Config(mbz.Region(*region), *clientID, *clientSecret)
		if err != nil {
			return err
		}
		token, err := config.Token(cmd.Context())
		if err != nil {
			return err
		}
		if err := writeCredentialsFile(&credentialsFile{
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

func newLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from the Mercedes-Benz Management API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := removeCredentialsFile(); err != nil {
				return err
			}
			cmd.Println("Logged out.")
			return nil
		},
	}
}

// credentialsFile storing authentication credentials for the CLI.
type credentialsFile struct {
	// Region is the region of the credentials.
	Region mbz.Region `json:"region"`
	// Credentials is the current stored client credentials.
	Credentials oauth2.Token `json:"clientCredentials"`
}

func resolveCredentialsFilepath() (string, error) {
	return xdg.ConfigFile("mbz-go/auth.json")
}

// readCredentialsFile reads the currently stored [credentialsFile].
func readCredentialsFile() (*credentialsFile, error) {
	credentialsFilepath, err := resolveCredentialsFilepath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(credentialsFilepath); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(credentialsFilepath)
	if err != nil {
		return nil, err
	}
	var cf credentialsFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, err
	}
	return &cf, nil
}

// writeCredentialsFile writes the stored [credentialsFile].
func writeCredentialsFile(cf *credentialsFile) error {
	credentialsFilepath, err := resolveCredentialsFilepath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credentialsFilepath, data, 0o600)
}

// removeCredentialsFile removes the stored [credentialsFile].
func removeCredentialsFile() error {
	credentialsFilepath, err := resolveCredentialsFilepath()
	if err != nil {
		return err
	}
	return os.RemoveAll(credentialsFilepath)
}
