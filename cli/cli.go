package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Store reads and writes JSON-serializable data.
type Store interface {
	Read(target any) error
	Write(data any) error
	Clear() error
}

// FleetCredentials holds OAuth2 + Kafka credentials for the Mercedes-Benz Fleet API.
type FleetCredentials struct {
	ClientID           string `json:"clientId"`
	ClientSecret       string `json:"clientSecret"`
	Region             string `json:"region"`
	KafkaConsumerGroup string `json:"kafkaConsumerGroup"`
	KafkaInputTopic    string `json:"kafkaInputTopic"`
}

// VehicleSpecCredentials holds API key credentials for the Mercedes-Benz Vehicle Specification API.
type VehicleSpecCredentials struct {
	APIKey string `json:"apiKey"`
}

// Option configures the CLI command tree.
type Option func(*config)

type config struct {
	fleetCredentialStore       Store
	vehicleSpecCredentialStore Store
	tokenStore                 Store
	httpClient                 *http.Client

	fleetCredentialProvider       func() (*FleetCredentials, error)
	vehicleSpecCredentialProvider func() (*VehicleSpecCredentials, error)
}

// WithFleetCredentialStore sets the credential store for Mercedes-Benz Fleet API (OAuth2 + Kafka).
func WithFleetCredentialStore(s Store) Option {
	return func(c *config) { c.fleetCredentialStore = s }
}

// WithVehicleSpecCredentialStore sets the credential store for Mercedes-Benz Vehicle Specification API.
func WithVehicleSpecCredentialStore(s Store) Option {
	return func(c *config) { c.vehicleSpecCredentialStore = s }
}

// WithTokenStore sets the token store.
func WithTokenStore(s Store) Option {
	return func(c *config) { c.tokenStore = s }
}

// WithHTTPClient sets a custom [http.Client] for the SDK client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *config) { c.httpClient = httpClient }
}

// WithFleetCredentialProvider sets a provider function for fleet credentials.
// When set, the provider is called instead of reading from the credential store.
func WithFleetCredentialProvider(fn func() (*FleetCredentials, error)) Option {
	return func(c *config) { c.fleetCredentialProvider = fn }
}

// WithVehicleSpecCredentialProvider sets a provider function for vehicle spec credentials.
// When set, the provider is called instead of reading from the credential store.
func WithVehicleSpecCredentialProvider(fn func() (*VehicleSpecCredentials, error)) Option {
	return func(c *config) { c.vehicleSpecCredentialProvider = fn }
}

// FileStore is a file-backed store that uses protojson for proto messages
// and encoding/json for other types.
type FileStore struct {
	path string
}

// NewFileStore creates a new file-backed store at the given path.
func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

// Read unmarshals the file contents into target.
// If target implements [proto.Message], protojson is used for unmarshaling.
func (s *FileStore) Read(target any) error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("read store: %w", err)
	}
	if msg, ok := target.(proto.Message); ok {
		if err := protojson.Unmarshal(data, msg); err != nil {
			return fmt.Errorf("unmarshal store: %w", err)
		}
		return nil
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("unmarshal store: %w", err)
	}
	return nil
}

// Write marshals data and writes it to the file.
// If data implements [proto.Message], protojson is used for marshaling.
func (s *FileStore) Write(data any) error {
	var bytes []byte
	if msg, ok := data.(proto.Message); ok {
		var err error
		bytes, err = protojson.MarshalOptions{
			Multiline: true,
			Indent:    "  ",
		}.Marshal(msg)
		if err != nil {
			return fmt.Errorf("marshal store: %w", err)
		}
	} else {
		var err error
		bytes, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal store: %w", err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create store dir: %w", err)
	}
	return os.WriteFile(s.path, bytes, 0o600)
}

// Clear removes the file.
func (s *FileStore) Clear() error {
	err := os.Remove(s.path)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}
