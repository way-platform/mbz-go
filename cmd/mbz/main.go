package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/oauth"
	"github.com/way-platform/mbz-go"
	"github.com/way-platform/mbz-go/api/vehiclesv1"
	"github.com/way-platform/mbz-go/cmd/mbz/internal/auth"
	"google.golang.org/protobuf/encoding/protojson"
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
	cmd.AddCommand(newGetVehicleServicesCommand())
	cmd.AddCommand(newPostVehicleServicesCommand())
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
	cmd.AddGroup(&cobra.Group{
		ID:    "utils",
		Title: "Utils",
	})
	cmd.SetHelpCommandGroupID("utils")
	cmd.SetCompletionCommandGroupID("utils")
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

func newGetVehicleServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get-vehicle-services",
		Short:   "Get vehicle services",
		GroupID: "vehicles",
		Args:    cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.GetVehicleServices(cmd.Context(), &mbz.GetVehicleServicesRequest{
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

func newPostVehicleServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "post-vehicle-services",
		Short:   "Post vehicle services",
		GroupID: "vehicles",
		Args:    cobra.MinimumNArgs(3),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		services := make([]mbz.VehicleServices, 0)
		for i := 1; i < len(args); i += 2 {
			if strings.ToLower(args[i+1]) != "active" && strings.ToLower(args[i+1]) != "inactive" {
				return fmt.Errorf("invalid desired status: %s", args[i+1])
			}
			services = append(services, mbz.VehicleServices{
				ServiceID:     args[i],
				DesiredStatus: mbz.DesiredStatus(strings.ToUpper(args[i+1])),
			})
		}
		response, err := client.PostVehicleServices(cmd.Context(), &mbz.PostVehicleServicesRequest{
			DesiredServiceStatusInput: []mbz.DesiredServiceStatusInput{
				{
					VIN:      args[0],
					Services: services,
				},
			},
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
	_ = cmd.MarkFlagRequired("topic")
	consumerGroup := cmd.Flags().String("consumer-group", "", "Consumer group")
	_ = cmd.MarkFlagRequired("consumer-group")
	enableDebug := cmd.Flags().Bool("debug", false, "Enable debug logging")
	format := cmd.Flags().String("format", "json", "Format to use for output")
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
		opts := []kgo.Opt{
			kgo.DialTLS(),
			kgo.SeedBrokers(bootstrapServer),
			kgo.ConsumerGroup(*consumerGroup),
			kgo.ConsumeTopics(*topic),
			kgo.SASL(oauth.Oauth(func(ctx context.Context) (oauth.Auth, error) {
				return oauth.Auth{
					Token: authFile.Credentials.AccessToken,
				}, nil
			})),
		}
		if *enableDebug {
			opts = append(opts, kgo.WithLogger(&logger{sl: slog.Default()}))
		}
		client, err := kgo.NewClient(opts...)
		if err != nil {
			return fmt.Errorf("failed to create Kafka client: %w", err)
		}
		defer client.Close()
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
		defer cancel()
		for {
			fetches := client.PollFetches(ctx)
			if fetches.IsClientClosed() || ctx.Err() != nil {
				break
			}
			var errs []error
			fetches.EachError(func(topic string, partition int32, err error) {
				errs = append(errs, err)
			})
			if len(errs) > 0 {
				return fmt.Errorf("errors fetching records: %w", errors.Join(errs...))
			}
			it := fetches.RecordIter()
			for !it.Done() {
				record := it.Next()
				switch *format {
				case "json":
					var msg mbz.PushMessage
					if err := json.Unmarshal(record.Value, &msg); err != nil {
						return fmt.Errorf("failed to unmarshal message: %w", err)
					}
					data, err := json.MarshalIndent(msg, "", "  ")
					if err != nil {
						return fmt.Errorf("failed to marshal message: %w", err)
					}
					fmt.Println(string(data))
				case "proto":
					var msg mbz.PushMessage
					if err := json.Unmarshal(record.Value, &msg); err != nil {
						return fmt.Errorf("failed to unmarshal message: %w", err)
					}
					protoMsg, err := msg.AsProto()
					if err != nil {
						return fmt.Errorf("failed to convert message to protobuf: %w", err)
					}
					fmt.Println(protojson.Format(protoMsg))
				case "raw":
					fmt.Println(string(record.Value))
				}
			}
			slog.Debug("fetched records", "count", fetches.NumRecords())
		}
		return nil
	}
	return cmd
}

func printJSON(msg any) {
	data, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		slog.Error("failed to marshal JSON", "error", err)
	}
	fmt.Println(string(data))
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
