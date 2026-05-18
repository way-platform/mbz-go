package mbz

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/way-platform/mbz-go/api/vehiclesv1"
	fleetv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/fleet/v1"
	"golang.org/x/oauth2"
)

func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	client := &Client{
		baseURL: srv.URL + "/api",
		config: clientConfig{
			tokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"}),
			retryCount:  0,
			timeout:     0,
		},
	}
	return client
}

func TestListVehicles(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/vehicles" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", got)
		}
		vehicles := []vehiclesv1.Vehicle{
			{VIN: "WDB1234567890001"},
			{VIN: "WDB1234567890002"},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(vehicles)
	}))
	t.Cleanup(srv.Close)

	client := newTestClient(t, srv)
	resp, err := client.ListVehicles(context.Background(), &fleetv1.ListVehiclesRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(resp.GetVehicles()); got != 2 {
		t.Fatalf("expected 2 vehicles, got %d", got)
	}
	if got := resp.GetVehicles()[0].GetVin(); got != "WDB1234567890001" {
		t.Errorf("expected WDB1234567890001, got %s", got)
	}
	if got := resp.GetVehicles()[1].GetVin(); got != "WDB1234567890002" {
		t.Errorf("expected WDB1234567890002, got %s", got)
	}
}

func TestAssignVehicles(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/vehicles" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body []vehiclesv1.Vehicle
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body) != 2 {
			t.Fatalf("expected 2 vehicles in body, got %d", len(body))
		}
		if body[0].VIN != "VIN1" || body[1].VIN != "VIN2" {
			t.Errorf("unexpected VINs: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(srv.Close)

	client := newTestClient(t, srv)
	req := &fleetv1.AssignVehiclesRequest{}
	req.SetVins([]string{"VIN1", "VIN2"})
	_, err := client.AssignVehicles(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteVehicles(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/vehicles" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body []vehiclesv1.Vehicle
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body) != 1 || body[0].VIN != "VIN_TO_DELETE" {
			t.Errorf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	client := newTestClient(t, srv)
	req := &fleetv1.DeleteVehiclesRequest{}
	req.SetVins([]string{"VIN_TO_DELETE"})
	_, err := client.DeleteVehicles(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResponseError_ConnectCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantCode   connect.Code
	}{
		{"bad request", http.StatusBadRequest, connect.CodeInvalidArgument},
		{"unauthorized", http.StatusUnauthorized, connect.CodeUnauthenticated},
		{"forbidden", http.StatusForbidden, connect.CodePermissionDenied},
		{"not found", http.StatusNotFound, connect.CodeNotFound},
		{"conflict", http.StatusConflict, connect.CodeAlreadyExists},
		{"too many requests", http.StatusTooManyRequests, connect.CodeResourceExhausted},
		{"internal server error", http.StatusInternalServerError, connect.CodeInternal},
		{"service unavailable", http.StatusServiceUnavailable, connect.CodeUnavailable},
		{"gateway timeout", http.StatusGatewayTimeout, connect.CodeDeadlineExceeded},
		{"unknown status", http.StatusTeapot, connect.CodeUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte("error body"))
			}))
			t.Cleanup(srv.Close)

			client := newTestClient(t, srv)
			_, err := client.ListVehicles(context.Background(), &fleetv1.ListVehiclesRequest{})
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var connectErr *connect.Error
			if !errors.As(err, &connectErr) {
				t.Fatalf("expected connect.Error, got %T: %v", err, err)
			}
			if connectErr.Code() != tt.wantCode {
				t.Errorf("expected code %v, got %v", tt.wantCode, connectErr.Code())
			}
		})
	}
}
