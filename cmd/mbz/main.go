package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
	"github.com/way-platform/mbz-go"
	"github.com/way-platform/mbz-go/api/vehiclesv1"
	"github.com/way-platform/mbz-go/cmd/mbz/internal/auth"
)

func main() {
	if err := fang.Execute(
		context.Background(),
		newRootCommand(),
		fang.WithColorSchemeFunc(fang.AnsiColorScheme),
	); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mbz",
		Short: "Mercedes-Benz Management API CLI",
	}
	cmd.AddGroup(&cobra.Group{
		ID:    "auth",
		Title: "Authentication",
	})
	authCmd := auth.NewCommand()
	authCmd.GroupID = "auth"
	cmd.AddCommand(authCmd)
	cmd.AddGroup(&cobra.Group{
		ID:    "vehicles",
		Title: "Vehicles",
	})
	cmd.AddCommand(newListVehiclesCommand())
	cmd.AddCommand(newAssignVehiclesCommand())
	cmd.AddCommand(newDeleteVehiclesCommand())
	cmd.AddCommand(newGetVehicleCompatibilityCommand())
	cmd.AddCommand(newEnableDeltaPushCommand())
	cmd.AddCommand(newDisableDeltaPushCommand())
	cmd.AddGroup(&cobra.Group{
		ID:    "services",
		Title: "Services",
	})
	cmd.AddCommand(newListServicesCommand())
	return cmd
}

func newListVehiclesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-vehicles",
		Short:   "List vehicles",
		GroupID: "vehicles",
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
		printJSON(response)
		return nil
	}
	return cmd
}

func newAssignVehiclesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assign-vehicles",
		Short:   "Assign vehicles",
		GroupID: "vehicles",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.AssignVehicles(cmd.Context(), &mbz.AssignVehiclesRequest{
			VINs: args,
		})
		if err != nil {
			return err
		}
		printJSON(response)
		return nil
	}
	return cmd
}

func newDeleteVehiclesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete-vehicles",
		Short:   "Delete vehicles",
		GroupID: "vehicles",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.DeleteVehicles(cmd.Context(), &mbz.DeleteVehiclesRequest{
			VINs: args,
		})
		if err != nil {
			return err
		}
		printJSON(response)
		return nil
	}
	return cmd
}

func newListServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-services",
		Short:   "List services",
		GroupID: "services",
	}
	details := cmd.Flags().Bool("details", false, "Include service details")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.ListServices(cmd.Context(), &mbz.ListServicesRequest{
			Details: *details,
		})
		if err != nil {
			return err
		}
		printJSON(response)
		return nil
	}
	return cmd
}

func newGetVehicleCompatibilityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get-vehicle-compatibility",
		Short:   "Get vehicle compatibility",
		GroupID: "vehicles",
		Args:    cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.GetVehicleCompatibility(cmd.Context(), &mbz.GetVehicleCompatibilityRequest{
			VIN: args[0],
		})
		if err != nil {
			return err
		}
		printJSON(response)
		return nil
	}
	return cmd
}

func newEnableDeltaPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "enable-delta-push",
		Short:   "Enable delta push",
		GroupID: "vehicles",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		request := mbz.PatchVehiclesRequest{
			Vehicles: make([]vehiclesv1.Vehicle, 0, len(args)),
		}
		for _, vin := range args {
			request.Vehicles = append(request.Vehicles, vehiclesv1.Vehicle{
				VIN:       vin,
				DeltaPush: ptr(true),
			})
		}
		response, err := client.PatchVehicles(cmd.Context(), &request)
		if err != nil {
			return err
		}
		printJSON(response)
		return nil
	}
	return cmd
}

func newDisableDeltaPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable-delta-push",
		Short:   "Disable delta push",
		GroupID: "vehicles",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		request := mbz.PatchVehiclesRequest{
			Vehicles: make([]vehiclesv1.Vehicle, 0, len(args)),
		}
		for _, vin := range args {
			request.Vehicles = append(request.Vehicles, vehiclesv1.Vehicle{
				VIN:       vin,
				DeltaPush: ptr(false),
			})
		}
		response, err := client.PatchVehicles(cmd.Context(), &request)
		if err != nil {
			return err
		}
		printJSON(response)
		return nil
	}
	return cmd
}

func printJSON(msg any) error {
	data, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func ptr[T any](v T) *T {
	return &v
}
