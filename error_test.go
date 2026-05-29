package mbz

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
	fleetv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/fleet/v1"
)

func TestResponseError_vehicleNotFound(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{
			"type": "/vehicle/not-found",
			"title": "Not found",
			"detail": "Server can not find the requested resource",
			"instance": "about:blank",
			"statusCode": 404
		}`)
	}))
	t.Cleanup(srv.Close)
	client := newTestClient(t, srv)
	_, err := client.ListVehicles(t.Context(), &fleetv1.ListVehiclesRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected connect error, got %T", err)
	}
	if connectErr.Code() != connect.CodeNotFound {
		t.Errorf("code = %v, want %v", connectErr.Code(), connect.CodeNotFound)
	}
	detail := firstDetail[*mbzv1.ProblemDetail](t, connectErr)
	if detail.GetType() != mbzv1.ProblemDetail_TYPE_VEHICLE_NOT_FOUND {
		t.Errorf("type = %v, want TYPE_VEHICLE_NOT_FOUND", detail.GetType())
	}
	if detail.GetTypeUri() != "/vehicle/not-found" {
		t.Errorf("type_uri = %q", detail.GetTypeUri())
	}
	if detail.GetTitle() != "Not found" {
		t.Errorf("title = %q", detail.GetTitle())
	}
	if detail.GetDetail() != "Server can not find the requested resource" {
		t.Errorf("detail = %q", detail.GetDetail())
	}
	if detail.GetStatus() != 404 {
		t.Errorf("status_code = %d", detail.GetStatus())
	}
}

func TestResponseError_validationError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{
			"type": "/error/validation-error",
			"title": "Invalid JSON Input",
			"detail": "Unexpected character (',' (code 44)): was expecting double-quote to start field name",
			"instance": "about:blank",
			"statusCode": 400
		}`)
	}))
	t.Cleanup(srv.Close)
	client := newTestClient(t, srv)
	_, err := client.ListVehicles(t.Context(), &fleetv1.ListVehiclesRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected connect error, got %T", err)
	}
	if connectErr.Code() != connect.CodeInvalidArgument {
		t.Errorf("code = %v, want %v", connectErr.Code(), connect.CodeInvalidArgument)
	}
	detail := firstDetail[*mbzv1.ProblemDetail](t, connectErr)
	if detail.GetType() != mbzv1.ProblemDetail_TYPE_VALIDATION_ERROR {
		t.Errorf("type = %v, want TYPE_VALIDATION_ERROR", detail.GetType())
	}
	if detail.GetTypeUri() != "/error/validation-error" {
		t.Errorf("type_uri = %q", detail.GetTypeUri())
	}
}

func TestResponseError_unrecognizedType(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{
			"type": "/error/some-new-error",
			"title": "Forbidden",
			"detail": "Access denied",
			"statusCode": 403
		}`)
	}))
	t.Cleanup(srv.Close)
	client := newTestClient(t, srv)
	_, err := client.ListVehicles(t.Context(), &fleetv1.ListVehiclesRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected connect error, got %T", err)
	}
	detail := firstDetail[*mbzv1.ProblemDetail](t, connectErr)
	if detail.GetType() != mbzv1.ProblemDetail_TYPE_UNRECOGNIZED {
		t.Errorf("type = %v, want TYPE_UNRECOGNIZED", detail.GetType())
	}
	if detail.GetTypeUri() != "/error/some-new-error" {
		t.Errorf("type_uri = %q, want /error/some-new-error", detail.GetTypeUri())
	}
}

func TestResponseError_nonJSONBody(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "Internal Server Error")
	}))
	t.Cleanup(srv.Close)
	client := newTestClient(t, srv)
	_, err := client.ListVehicles(t.Context(), &fleetv1.ListVehiclesRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected connect error, got %T", err)
	}
	if connectErr.Code() != connect.CodeInternal {
		t.Errorf("code = %v, want %v", connectErr.Code(), connect.CodeInternal)
	}
	if len(connectErr.Details()) != 0 {
		t.Errorf("expected no details for non-JSON body, got %d", len(connectErr.Details()))
	}
}

func TestResponseError_emptyAboutBlankType(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{
			"type": "about:blank",
			"title": "Not found",
			"detail": "Vehicle do not exist or are not in your ownership",
			"statusCode": 404
		}`)
	}))
	t.Cleanup(srv.Close)
	client := newTestClient(t, srv)
	_, err := client.ListVehicles(context.Background(), &fleetv1.ListVehiclesRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected connect error, got %T", err)
	}
	detail := firstDetail[*mbzv1.ProblemDetail](t, connectErr)
	if detail.GetType() != mbzv1.ProblemDetail_TYPE_UNRECOGNIZED {
		t.Errorf("type = %v, want TYPE_UNRECOGNIZED", detail.GetType())
	}
	if detail.GetTypeUri() != "about:blank" {
		t.Errorf("type_uri = %q", detail.GetTypeUri())
	}
	if detail.GetDetail() != "Vehicle do not exist or are not in your ownership" {
		t.Errorf("detail = %q", detail.GetDetail())
	}
}

func firstDetail[T any](t *testing.T, err *connect.Error) T {
	t.Helper()
	for _, d := range err.Details() {
		value, valErr := d.Value()
		if valErr != nil {
			continue
		}
		if typed, ok := value.(T); ok {
			return typed
		}
	}
	var zero T
	t.Fatalf("no detail of type %T found in error details", zero)
	return zero
}
