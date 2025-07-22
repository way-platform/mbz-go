package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"log/slog"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/oauth"
	"github.com/way-platform/mbz-go"
	"github.com/way-platform/mbz-go/api/vehiclesv1"
	"github.com/way-platform/mbz-go/cmd/mbz/internal/auth"
)

func main() {
	if err := fang.Execute(
		context.Background(),
		newRootCommand(),
		fang.WithColorSchemeFunc(func(c lipgloss.LightDarkFunc) fang.ColorScheme {
			base := c(lipgloss.Black, lipgloss.White)
			baseInverted := c(lipgloss.White, lipgloss.Black)
			return fang.ColorScheme{
				Base:         base,
				Title:        base,
				Description:  base,
				Comment:      base,
				Flag:         base,
				FlagDefault:  base,
				Command:      base,
				QuotedString: base,
				Argument:     base,
				Help:         base,
				Dash:         base,
				ErrorHeader:  [2]color.Color{baseInverted, base},
				ErrorDetails: base,
			}
		}),
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
		ID:    "vehicle-signals",
		Title: "Vehicle Signals",
	})
	cmd.AddCommand(newConsumeVehicleSignalsCommand())
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

func newConsumeVehicleSignalsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "consume-vehicle-signals",
		Short:   "Consume vehicle signals",
		GroupID: "vehicle-signals",
	}
	topic := cmd.Flags().String("topic", "", "Topic")
	cmd.MarkFlagRequired("topic")
	consumerGroup := cmd.Flags().String("consumer-group", "", "Consumer group")
	cmd.MarkFlagRequired("consumer-group")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		authFile, err := auth.ReadFile()
		if err != nil {
			return err
		}
		var bootstrapServer string
		switch authFile.Region {
		case mbz.RegionECE:
			bootstrapServer = mbz.KafkaBootstrapServerECE
		case mbz.RegionAMAPNA:
			bootstrapServer = mbz.KafkaBootstrapServerAMAPNA
		default:
			return fmt.Errorf("invalid region: %s", authFile.Region)
		}
		client, err := kgo.NewClient(
			kgo.SeedBrokers(bootstrapServer),
			kgo.ConsumerGroup(*consumerGroup),
			kgo.ConsumeTopics(*topic),
			kgo.SASL(oauth.Oauth(func(ctx context.Context) (oauth.Auth, error) {
				return oauth.Auth{
					Token: authFile.Credentials.AccessToken,
				}, nil
			})),
			kgo.WithLogger(&logger{sl: slog.Default()}),
		)
		if err != nil {
			log.Fatalf("Failed to create Kafka client: %v", err)
		}
		defer client.Close()
		for {
			slog.Info("polling Kafka")
			fetches := client.PollFetches(cmd.Context())
			if err := fetches.Err(); err != nil {
				return err
			}
			slog.Info("fetched records", "count", fetches.NumRecords())
			it := fetches.RecordIter()
			for !it.Done() {
				record := it.Next()
				fmt.Println(string(record.Value))
			}
		}
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

type logger struct {
	sl *slog.Logger
}

func (l *logger) Level() kgo.LogLevel {
	ctx := context.Background()
	switch {
	case l.sl.Enabled(ctx, slog.LevelDebug):
		return kgo.LogLevelDebug
	case l.sl.Enabled(ctx, slog.LevelInfo):
		return kgo.LogLevelInfo
	case l.sl.Enabled(ctx, slog.LevelWarn):
		return kgo.LogLevelWarn
	case l.sl.Enabled(ctx, slog.LevelError):
		return kgo.LogLevelError
	default:
		return kgo.LogLevelNone
	}
}

func (l *logger) Log(level kgo.LogLevel, msg string, keyvals ...any) {
	l.sl.Log(context.Background(), kgoToSlogLevel(level), msg, keyvals...)
}

func kgoToSlogLevel(level kgo.LogLevel) slog.Level {
	switch level {
	case kgo.LogLevelError:
		return slog.LevelError
	case kgo.LogLevelWarn:
		return slog.LevelWarn
	case kgo.LogLevelInfo:
		return slog.LevelInfo
	case kgo.LogLevelDebug:
		return slog.LevelDebug
	default:
		// Using the default level for slog
		return slog.LevelInfo
	}
}
