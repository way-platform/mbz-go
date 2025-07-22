package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/way-platform/mbz-go"
	"github.com/way-platform/mbz-go/cmd/mbz/internal/auth"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mbz",
		Short: "Mercedes-Benz Management API CLI",
	}
	cmd.AddCommand(auth.NewCommand())
	cmd.AddCommand(newListVehiclesCommand())
	return cmd
}

func newListVehiclesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-vehicles",
		Short: "List vehicles",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.ListVehicles(cmd.Context(), &mbz.ListVehiclesRequest{})
		if err != nil {
			return err
		}
		printJSON(cmd, response)
		return nil
	}
	return cmd
}

func printJSON(cmd *cobra.Command, msg any) error {
	data, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
