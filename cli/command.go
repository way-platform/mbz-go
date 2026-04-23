package cli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/oauth"
	"github.com/way-platform/mbz-go"
	fleetv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/fleet/v1"
	vehiclespecv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/vehiclespec/v1"
	"golang.org/x/oauth2"
	"golang.org/x/term"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewCommand builds the full CLI command tree for the Mercedes-Benz SDK.
func NewCommand(opts ...Option) *cobra.Command {
	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}
	cmd := &cobra.Command{
		Use:   "mbz",
		Short: "Mercedes-Benz API CLI",
	}
	cmd.AddGroup(&cobra.Group{ID: "vehicles", Title: "Vehicles"})
	cmd.AddCommand(newListVehiclesCommand(&cfg))
	cmd.AddCommand(newAssignVehiclesCommand(&cfg))
	cmd.AddCommand(newDeleteVehiclesCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "services", Title: "Data Services"})
	cmd.AddCommand(newListServicesCommand(&cfg))
	cmd.AddCommand(newGetVehicleCompatibilityCommand(&cfg))
	cmd.AddCommand(newGetVehicleServicesCommand(&cfg))
	cmd.AddCommand(newActivateVehicleServicesCommand(&cfg))
	cmd.AddCommand(newDeactivateVehicleServicesCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "delta-push", Title: "Delta Push"})
	cmd.AddCommand(newEnableDeltaPushCommand(&cfg))
	cmd.AddCommand(newDisableDeltaPushCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "vehicle-specifications", Title: "Vehicle Specifications"})
	cmd.AddCommand(newGetVehicleSpecificationCommand(&cfg))
	cmd.AddCommand(newGetVehicleImagesCommand(&cfg))
	cmd.AddCommand(newGetVehicleImageCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "kafka", Title: "Kafka"})
	cmd.AddCommand(newConsumeVehicleSignalsCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "auth", Title: "Authentication"})
	cmd.AddCommand(newAuthCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "utils", Title: "Utils"})
	cmd.SetHelpCommandGroupID("utils")
	cmd.SetCompletionCommandGroupID("utils")
	return cmd
}

// Auth commands.

func newAuthCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth",
		Short:   "Authenticate to the Mercedes-Benz API",
		GroupID: "auth",
	}
	cmd.AddCommand(newLoginCommand(cfg))
	cmd.AddCommand(newLogoutCommand(cfg))
	return cmd
}

func newLoginCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the Mercedes-Benz API",
	}
	cmd.AddCommand(newLoginFleetCommand(cfg))
	cmd.AddCommand(newLoginVehicleSpecCommand(cfg))
	return cmd
}

func newLoginFleetCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fleet",
		Short: "Login to the Mercedes-Benz Fleet API (OAuth2)",
	}
	region := cmd.Flags().String("region", "", "Region (e.g. ECE, AMAP/NA)")
	clientID := cmd.Flags().String("client-id", "", "OAuth2 client ID")
	clientSecret := cmd.Flags().String("client-secret", "", "OAuth2 client secret")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		creds := &FleetCredentials{}
		if cfg.fleetCredentialStore != nil {
			if loaded, err := cfg.fleetCredentialStore.Load(); err == nil {
				creds = loaded
			} else if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("read credentials: %w", err)
			}
		}
		// Override with flags.
		if *region != "" {
			creds.Region = *region
		}
		if *clientID != "" {
			creds.ClientID = *clientID
		}
		if *clientSecret != "" {
			creds.ClientSecret = *clientSecret
		}
		// Default region.
		if creds.Region == "" {
			creds.Region = string(mbz.RegionECE)
		}
		// Prompt for missing fields.
		if creds.ClientID == "" {
			val, err := promptSecret(cmd, "Enter OAuth2 client ID: ")
			if err != nil {
				return err
			}
			creds.ClientID = val
		}
		if creds.ClientSecret == "" {
			val, err := promptSecret(cmd, "Enter OAuth2 client secret: ")
			if err != nil {
				return err
			}
			creds.ClientSecret = val
		}
		// Persist credentials.
		if cfg.fleetCredentialStore != nil {
			if err := cfg.fleetCredentialStore.Save(creds); err != nil {
				return fmt.Errorf("write credentials: %w", err)
			}
		}
		// Run OAuth2 flow.
		oauth2Config, err := mbz.NewOAuth2Config(
			mbz.Region(creds.Region),
			creds.ClientID,
			creds.ClientSecret,
		)
		if err != nil {
			return err
		}
		token, err := oauth2Config.Token(cmd.Context())
		if err != nil {
			return err
		}
		if cfg.tokenStore != nil {
			if err := cfg.tokenStore.Save(token); err != nil {
				return fmt.Errorf("write token: %w", err)
			}
		}
		cmd.Printf("Logged in to %s.\n", creds.Region)
		return nil
	}
	return cmd
}

func newLoginVehicleSpecCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vehicle-spec",
		Short: "Login to the Mercedes-Benz Vehicle Specification API (API key)",
	}
	apiKey := cmd.Flags().String("api-key", "", "API key")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		creds := &VehicleSpecCredentials{}
		if cfg.vehicleSpecCredentialStore != nil {
			if loaded, err := cfg.vehicleSpecCredentialStore.Load(); err == nil {
				creds = loaded
			} else if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("read credentials: %w", err)
			}
		}
		if *apiKey != "" {
			creds.APIKey = *apiKey
		}
		if creds.APIKey == "" {
			val, err := promptSecret(cmd, "Enter API key: ")
			if err != nil {
				return err
			}
			creds.APIKey = val
		}
		if cfg.vehicleSpecCredentialStore != nil {
			if err := cfg.vehicleSpecCredentialStore.Save(creds); err != nil {
				return fmt.Errorf("write credentials: %w", err)
			}
		}
		cmd.Println("Logged in to Vehicle Specification API.")
		return nil
	}
	return cmd
}

func newLogoutCommand(cfg *config) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from the Mercedes-Benz API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cfg.tokenStore != nil {
				if err := cfg.tokenStore.Clear(); err != nil {
					return fmt.Errorf("clear token: %w", err)
				}
			}
			if cfg.fleetCredentialStore != nil {
				if err := cfg.fleetCredentialStore.Clear(); err != nil {
					return fmt.Errorf("clear fleet credentials: %w", err)
				}
			}
			if cfg.vehicleSpecCredentialStore != nil {
				if err := cfg.vehicleSpecCredentialStore.Clear(); err != nil {
					return fmt.Errorf("clear vehicle-spec credentials: %w", err)
				}
			}
			cmd.Println("Logged out.")
			return nil
		},
	}
}

// Vehicle commands.

func newListVehiclesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vehicles",
		Short:   "List vehicles",
		GroupID: "vehicles",
	}
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		response, err := client.ListVehicles(cmd.Context(), &fleetv1.ListVehiclesRequest{})
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

func newAssignVehiclesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assign-vehicles <vin...>",
		Short:   "Assign vehicles",
		GroupID: "vehicles",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		req := &fleetv1.AssignVehiclesRequest{}
		req.SetVins(args)
		response, err := client.AssignVehicles(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

func newDeleteVehiclesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete-vehicles <vin...>",
		Short:   "Delete vehicles",
		GroupID: "vehicles",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		req := &fleetv1.DeleteVehiclesRequest{}
		req.SetVins(args)
		response, err := client.DeleteVehicles(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

// Service commands.

func newListServicesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "services",
		Short:   "List all available data services",
		GroupID: "services",
	}
	details := cmd.Flags().Bool("details", false, "Include service details")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		req := &fleetv1.ListServicesRequest{}
		req.SetDetails(*details)
		response, err := client.ListServices(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

func newGetVehicleCompatibilityCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vehicle-compatibility <vin>",
		Short:   "Get data service compatibility for a specific vehicle",
		GroupID: "services",
		Args:    cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		req := &fleetv1.GetVehicleCompatibilityRequest{}
		req.SetVin(args[0])
		response, err := client.GetVehicleCompatibility(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response.GetVehicleCompatibility()))
		return nil
	}
	return cmd
}

func newGetVehicleServicesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vehicle-services <vin>",
		Short:   "Get data services for a specific vehicle",
		GroupID: "services",
		Args:    cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		req := &fleetv1.GetVehicleServicesRequest{}
		req.SetVin(args[0])
		response, err := client.GetVehicleServices(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

func newActivateVehicleServicesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "activate-vehicle-services <vin> <service-id...>",
		Short:   "Activate vehicle services",
		GroupID: "services",
		Args:    cobra.MinimumNArgs(2),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		vin := args[0]
		services := make([]*fleetv1.VehicleServiceDesiredStatus, 0, len(args)-1)
		for i := 1; i < len(args); i++ {
			svc := &fleetv1.VehicleServiceDesiredStatus{}
			svc.SetServiceId(args[i])
			svc.SetDesiredStatus("ACTIVE")
			services = append(services, svc)
		}
		input := &fleetv1.VehicleServiceInput{}
		input.SetVin(vin)
		input.SetServices(services)
		req := &fleetv1.PostVehicleServicesRequest{}
		req.SetVehicleServiceInputs([]*fleetv1.VehicleServiceInput{input})
		response, err := client.PostVehicleServices(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

func newDeactivateVehicleServicesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deactivate-vehicle-services <vin> <service-id...>",
		Short:   "Deactivate vehicle services",
		GroupID: "services",
		Args:    cobra.MinimumNArgs(2),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		vin := args[0]
		services := make([]*fleetv1.VehicleServiceDesiredStatus, 0, len(args)-1)
		for i := 1; i < len(args); i++ {
			svc := &fleetv1.VehicleServiceDesiredStatus{}
			svc.SetServiceId(args[i])
			svc.SetDesiredStatus("INACTIVE")
			services = append(services, svc)
		}
		input := &fleetv1.VehicleServiceInput{}
		input.SetVin(vin)
		input.SetServices(services)
		req := &fleetv1.PostVehicleServicesRequest{}
		req.SetVehicleServiceInputs([]*fleetv1.VehicleServiceInput{input})
		response, err := client.PostVehicleServices(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

// Delta push commands.

func newEnableDeltaPushCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "enable-delta-push <vin...>",
		Short:   "Enable delta push",
		GroupID: "delta-push",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		vehicles := make([]*fleetv1.Vehicle, 0, len(args))
		for _, vin := range args {
			v := &fleetv1.Vehicle{}
			v.SetVin(vin)
			v.SetDeltaPush(true)
			vehicles = append(vehicles, v)
		}
		req := &fleetv1.PatchVehiclesRequest{}
		req.SetVehicles(vehicles)
		response, err := client.PatchVehicles(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

func newDisableDeltaPushCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable-delta-push <vin...>",
		Short:   "Disable delta push",
		GroupID: "delta-push",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newOAuth2Client(cmd, cfg)
		if err != nil {
			return err
		}
		vehicles := make([]*fleetv1.Vehicle, 0, len(args))
		for _, vin := range args {
			v := &fleetv1.Vehicle{}
			v.SetVin(vin)
			v.SetDeltaPush(false)
			vehicles = append(vehicles, v)
		}
		req := &fleetv1.PatchVehiclesRequest{}
		req.SetVehicles(vehicles)
		response, err := client.PatchVehicles(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response))
		return nil
	}
	return cmd
}

// Vehicle specification commands.

func newGetVehicleSpecificationCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vehicle-specification <vin>",
		Short:   "Get vehicle specification",
		GroupID: "vehicle-specifications",
		Args:    cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClientWithAPIKey(cmd, cfg)
		if err != nil {
			return err
		}
		req := &vehiclespecv1.GetVehicleSpecificationRequest{}
		req.SetVin(args[0])
		response, err := client.GetVehicleSpecification(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response.GetVehicleSpecification()))
		return nil
	}
	return cmd
}

func newGetVehicleImagesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vehicle-images <vin>",
		Short:   "Get vehicle image IDs and URLs",
		GroupID: "vehicle-specifications",
		Args:    cobra.ExactArgs(1),
	}
	background := cmd.Flags().
		Bool("background", false, "Include background in images (high detail with realistic reflections)")
	fileFormat := cmd.Flags().String("file-format", "webp", "Image file format (png, jpeg, webp)")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClientWithAPIKey(cmd, cfg)
		if err != nil {
			return err
		}
		req := &vehiclespecv1.GetVehicleImageIdsRequest{}
		req.SetVin(args[0])
		req.SetBackground(*background)
		req.SetFileFormat(*fileFormat)
		response, err := client.GetVehicleImageIds(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Println(protojson.Format(response.GetVehicleImages()))
		return nil
	}
	return cmd
}

func newGetVehicleImageCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "image <image-id>",
		Short:   "Download vehicle image by image ID",
		GroupID: "vehicle-specifications",
		Args:    cobra.ExactArgs(1),
	}
	outputFile := cmd.Flags().
		StringP("output", "o", "", "Output file path (default: <image-id>.<extension>)")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClientWithAPIKey(cmd, cfg)
		if err != nil {
			return err
		}
		imageID := args[0]
		req := &vehiclespecv1.GetImageRequest{}
		req.SetImageId(imageID)
		response, err := client.GetImage(cmd.Context(), req)
		if err != nil {
			return err
		}
		outputPath := *outputFile
		if outputPath == "" {
			ext := getExtensionFromContentType(response.GetContentType())
			outputPath = imageID + ext
		}
		if err := os.WriteFile(outputPath, response.GetData(), 0o644); err != nil {
			return fmt.Errorf("failed to write image to file: %w", err)
		}
		cmd.Printf("Image downloaded to %s (%s)\n", outputPath, response.GetContentType())
		return nil
	}
	return cmd
}

// Kafka commands.

func newConsumeVehicleSignalsCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "consume-vehicle-signals",
		Short:   "Consume vehicle signals from Kafka",
		GroupID: "kafka",
	}
	topicFlag := cmd.Flags().String("topic", "", "Topic (overrides credential value)")
	consumerGroupFlag := cmd.Flags().String("consumer-group", "", "Consumer group (overrides credential value)")
	enableDebug := cmd.Flags().Bool("debug", false, "Enable debug logging")
	format := cmd.Flags().String("format", "json", "Format to use for output")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		creds, err := resolveFleetCredentials(cfg)
		if err != nil {
			return err
		}
		var token *oauth2.Token
		if cfg.tokenStore != nil {
			token, err = cfg.tokenStore.Load()
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return fmt.Errorf("no credentials found, please login using `mbz auth login fleet`")
				}
				return fmt.Errorf("read token: %w", err)
			}
		}
		region, err := resolveOAuth2Region(creds, token)
		if err != nil {
			return err
		}
		// Resolve topic and consumer group from credentials, with flag overrides.
		topic := creds.KafkaInputTopic
		if *topicFlag != "" {
			topic = *topicFlag
		}
		if topic == "" {
			return fmt.Errorf("topic is required: set via --topic or store in fleet credentials")
		}
		consumerGroup := creds.KafkaConsumerGroup
		if *consumerGroupFlag != "" {
			consumerGroup = *consumerGroupFlag
		}
		if consumerGroup == "" {
			return fmt.Errorf(
				"consumer group is required: set via --consumer-group or store in fleet credentials",
			)
		}
		var bootstrapServer string
		switch region {
		case mbz.RegionECE:
			bootstrapServer = mbz.KafkaBootstrapServerECE
		case mbz.RegionAMAPNA:
			bootstrapServer = mbz.KafkaBootstrapServerAMAPNA
		default:
			return fmt.Errorf("invalid region: %s", creds.Region)
		}
		opts := []kgo.Opt{
			kgo.DialTLS(),
			kgo.SeedBrokers(bootstrapServer),
			kgo.ConsumerGroup(consumerGroup),
			kgo.ConsumeTopics(topic),
			kgo.SASL(oauth.Oauth(func(_ context.Context) (oauth.Auth, error) {
				var accessToken string
				if token != nil {
					accessToken = token.AccessToken
				}
				return oauth.Auth{
					Token: accessToken,
				}, nil
			})),
		}
		if *enableDebug {
			opts = append(opts, kgo.WithLogger(&kgoLogger{sl: slog.Default()}))
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
			fetches.EachError(func(_ string, _ int32, err error) {
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

// Client constructors.

func resolveFleetCredentials(cfg *config) (*FleetCredentials, error) {
	if cfg.fleetCredentialStore == nil {
		return nil, fmt.Errorf("no fleet credentials configured")
	}
	creds, err := cfg.fleetCredentialStore.Load()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("no credentials found, please login using `mbz auth login fleet`")
		}
		return nil, fmt.Errorf("read credentials: %w", err)
	}
	return creds, nil
}

func resolveVehicleSpecCredentials(cfg *config) (*VehicleSpecCredentials, error) {
	if cfg.vehicleSpecCredentialStore == nil {
		return nil, fmt.Errorf("no vehicle spec credentials configured")
	}
	creds, err := cfg.vehicleSpecCredentialStore.Load()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf(
				"no credentials found, please login using `mbz auth login vehicle-spec --api-key <key>`",
			)
		}
		return nil, fmt.Errorf("read credentials: %w", err)
	}
	return creds, nil
}

func newOAuth2Client(cmd *cobra.Command, cfg *config) (*mbz.Client, error) {
	creds, err := resolveFleetCredentials(cfg)
	if err != nil {
		return nil, err
	}
	var token *oauth2.Token
	if cfg.tokenStore != nil {
		token, err = cfg.tokenStore.Load()
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("no credentials found, please login using `mbz auth login fleet`")
			}
			return nil, fmt.Errorf("read token: %w", err)
		}
	}
	if token == nil || token.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("invalid token, please login using `mbz auth login fleet`")
	}
	region, err := resolveOAuth2Region(creds, token)
	if err != nil {
		return nil, err
	}
	opts := []mbz.ClientOption{
		mbz.WithRegion(region),
		mbz.WithOAuth2TokenSource(oauth2.StaticTokenSource(token)),
	}
	if cfg.httpClient != nil {
		opts = append(opts, mbz.WithHTTPClient(cfg.httpClient))
	}
	return mbz.NewClient(cmd.Context(), opts...)
}

func newClientWithAPIKey(cmd *cobra.Command, cfg *config) (*mbz.Client, error) {
	creds, err := resolveVehicleSpecCredentials(cfg)
	if err != nil {
		return nil, err
	}
	if creds.APIKey == "" {
		return nil, fmt.Errorf(
			"no API key found, please login using `mbz auth login vehicle-spec --api-key <key>`",
		)
	}
	opts := []mbz.ClientOption{
		mbz.WithAPIKey(creds.APIKey),
	}
	if cfg.httpClient != nil {
		opts = append(opts, mbz.WithHTTPClient(cfg.httpClient))
	}
	return mbz.NewClient(cmd.Context(), opts...)
}

// Helpers.

func promptSecret(cmd *cobra.Command, prompt string) (string, error) {
	cmd.Print(prompt)
	input, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	cmd.Println()
	return string(input), nil
}

func resolveOAuth2Region(creds *FleetCredentials, token *oauth2.Token) (mbz.Region, error) {
	if creds.Region != "" {
		return mbz.Region(creds.Region), nil
	}
	if token == nil {
		return "", fmt.Errorf("missing region and token")
	}
	return inferRegionFromAccessToken(token.AccessToken)
}

func inferRegionFromAccessToken(accessToken string) (mbz.Region, error) {
	if accessToken == "" {
		return "", fmt.Errorf("missing region and access token")
	}
	parts := strings.Split(accessToken, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("infer region from token: invalid jwt")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("infer region from token payload: %w", err)
	}
	var claims struct {
		Iss string `json:"iss"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("infer region from token claims: %w", err)
	}
	switch claims.Iss {
	case "https://ssoalpha.dvb.corpinter.net/v1":
		return mbz.RegionECE, nil
	case "https://ssoalpha.am.dvb.corpinter.net/v1":
		return mbz.RegionAMAPNA, nil
	default:
		return "", fmt.Errorf("infer region from token issuer: unknown issuer %q", claims.Iss)
	}
}

func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}

type kgoLogger struct {
	sl *slog.Logger
}

func (l *kgoLogger) Level() kgo.LogLevel {
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

func (l *kgoLogger) Log(level kgo.LogLevel, msg string, keyvals ...any) {
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
		return slog.LevelInfo
	}
}
