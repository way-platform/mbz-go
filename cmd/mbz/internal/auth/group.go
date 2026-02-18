package auth

import "github.com/spf13/cobra"

// NewGroup returns a new [cobra.Group] for CLI authentication.
func NewGroup() *cobra.Group {
	return &cobra.Group{
		ID:    "auth",
		Title: "Authentication",
	}
}
