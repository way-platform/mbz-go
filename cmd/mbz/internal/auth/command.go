package auth

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/way-platform/mbz-go"
	"golang.org/x/term"
)

// NewCommand returns a new [cobra.Command] for CLI authentication.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth",
		Short:   "Authenticate to the Mercedes-Benz API",
		GroupID: "auth",
	}
	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newLogoutCommand())
	return cmd
}

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the Mercedes-Benz API",
	}
	apiKey := cmd.Flags().String("api-key", "", "API key to use for authentication")
	region := cmd.Flags().String("region", string(mbz.RegionECE), "region to use for authentication (OAuth2 only)")
	clientID := cmd.Flags().String("client-id", "-", "client ID to use for authentication (OAuth2 only)")
	clientSecret := cmd.Flags().String("client-secret", "-", "client secret to use for authentication (OAuth2 only)")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if *apiKey == "-" {
			cmd.Print("Enter API key: ")
			input, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			*apiKey = string(input)
			cmd.Println()
		}
		// Otherwise, use OAuth2 authentication
		if *clientID == "-" && *apiKey == "" {
			cmd.Print("Enter OAuth2 client ID: ")
			input, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			*clientID = string(input)
			cmd.Println()
		}
		if *clientSecret == "-" && *apiKey == "" {
			cmd.Print("Enter OAuth2 client secret: ")
			input, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			*clientSecret = string(input)
			cmd.Println()
		}
		authFile := File{
			Region: mbz.Region(*region),
			APIKey: *apiKey,
		}
		if *clientSecret != "-" && *clientID != "-" {
			config, err := mbz.NewOAuth2Config(mbz.Region(*region), *clientID, *clientSecret)
			if err != nil {
				return err
			}
			token, err := config.Token(cmd.Context())
			if err != nil {
				return err
			}
			authFile.Token = token
		}
		if err := writeFile(&authFile); err != nil {
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
