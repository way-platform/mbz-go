package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// FleetCredentials holds Mercedes-Benz Fleet API credentials (OAuth2 + Kafka).
type FleetCredentials struct {
	ClientID           string `json:"clientId"`
	ClientSecret       string `json:"clientSecret"`
	Region             string `json:"region"`
	KafkaConsumerGroup string `json:"kafkaConsumerGroup"`
	KafkaInputTopic    string `json:"kafkaInputTopic"`
}

// VehicleSpecCredentials holds Mercedes-Benz Vehicle Specification API credentials.
type VehicleSpecCredentials struct {
	APIKey string `json:"apiKey"`
}

// FleetCredentialStore reads and writes fleet credentials.
type FleetCredentialStore interface {
	Load() (*FleetCredentials, error)
	Save(*FleetCredentials) error
	Clear() error
}

// VehicleSpecCredentialStore reads and writes vehicle specification credentials.
type VehicleSpecCredentialStore interface {
	Load() (*VehicleSpecCredentials, error)
	Save(*VehicleSpecCredentials) error
	Clear() error
}

// Store reads and writes JSON-serializable data.
type Store interface {
	Read(target any) error
	Write(data any) error
	Clear() error
}

// Option configures the CLI command tree.
type Option func(*config)

type config struct {
	fleetCredentialStore       FleetCredentialStore
	vehicleSpecCredentialStore VehicleSpecCredentialStore
	tokenStore                 Store
	httpClient                 *http.Client
}

// WithFleetCredentialStore sets the credential store for Mercedes-Benz Fleet API (OAuth2 + Kafka).
func WithFleetCredentialStore(s FleetCredentialStore) Option {
	return func(c *config) { c.fleetCredentialStore = s }
}

// WithVehicleSpecCredentialStore sets the credential store for Mercedes-Benz Vehicle Specification API.
func WithVehicleSpecCredentialStore(s VehicleSpecCredentialStore) Option {
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

// fileStore is a generic file-backed credential store using encoding/json.
type fileStore[T any] struct {
	path string
}

// Load reads and unmarshals credentials from the file.
func (s *fileStore[T]) Load() (*T, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("read store: %w", err)
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("unmarshal store: %w", err)
	}
	return &v, nil
}

// Save marshals and writes credentials to the file.
func (s *fileStore[T]) Save(v *T) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal store: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create store dir: %w", err)
	}
	return os.WriteFile(s.path, data, 0o600)
}

// Clear removes the file.
func (s *fileStore[T]) Clear() error {
	err := os.Remove(s.path)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}

// NewFleetCredentialFileStore creates a file-backed store for fleet credentials.
func NewFleetCredentialFileStore(path string) FleetCredentialStore {
	return &fileStore[FleetCredentials]{path: path}
}

// NewVehicleSpecCredentialFileStore creates a file-backed store for vehicle spec credentials.
func NewVehicleSpecCredentialFileStore(path string) VehicleSpecCredentialStore {
	return &fileStore[VehicleSpecCredentials]{path: path}
}

// FileStore is a file-backed store that uses encoding/json.
type FileStore struct {
	path string
}

// NewFileStore creates a new file-backed store at the given path.
func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

// Read unmarshals the file contents into target.
func (s *FileStore) Read(target any) error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("read store: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("unmarshal store: %w", err)
	}
	return nil
}

// Write marshals data and writes it to the file.
func (s *FileStore) Write(data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal store: %w", err)
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
