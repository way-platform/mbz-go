package mbz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"connectrpc.com/connect"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
	"google.golang.org/protobuf/proto"
)

// ResponseError represents a Mercedes-Benz API response error.
type ResponseError struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Body is the raw response body.
	Body []byte
	// Path is the API path that returned the error.
	Path string
}

func newResponseError(httpResponse *http.Response) error {
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		body = fmt.Appendf(nil, "failed to read response body: %s", err)
	}
	respErr := &ResponseError{
		StatusCode: httpResponse.StatusCode,
		Body:       body,
	}
	if httpResponse.Request != nil && httpResponse.Request.URL != nil {
		respErr.Path = httpResponse.Request.URL.Path
	}
	return respErr.connectError()
}

func (e *ResponseError) Error() string {
	if len(e.Body) > 0 {
		return fmt.Sprintf("http %d: %s", e.StatusCode, string(e.Body))
	}
	return fmt.Sprintf("http %d", e.StatusCode)
}

func (e *ResponseError) connectError() error {
	connErr := connect.NewError(httpStatusToConnectCode(e.StatusCode), e)
	if detail := e.errorDetail(); detail != nil {
		d, err := connect.NewErrorDetail(detail)
		if err == nil {
			connErr.AddDetail(d)
		}
	}
	return connErr
}

func (e *ResponseError) errorDetail() proto.Message {
	var raw problemDetailJSON
	if err := json.Unmarshal(e.Body, &raw); err != nil {
		return nil
	}
	if raw.Type == "" && raw.Title == "" && raw.Detail == "" {
		return nil
	}
	problemType := apiEnumValueToProto(raw.Type)
	b := mbzv1.ProblemDetail_builder{
		Type:     &problemType,
		Title:    optString(raw.Title),
		Detail:   optString(raw.Detail),
		Instance: optString(raw.Instance),
		Status:   optInt32(raw.status()),
	}
	if raw.Type != "" {
		b.TypeUri = &raw.Type
	}
	return b.Build()
}

type problemDetailJSON struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	Instance   string `json:"instance"`
	Status     int32  `json:"status"`
	StatusCode int32  `json:"statusCode"`
}

func (p *problemDetailJSON) status() int32 {
	if p.Status != 0 {
		return p.Status
	}
	return p.StatusCode
}

func apiEnumValueToProto(typeURI string) mbzv1.ProblemDetail_Type {
	if typeURI == "" {
		return mbzv1.ProblemDetail_TYPE_UNSPECIFIED
	}
	enumDesc := mbzv1.ProblemDetail_TYPE_UNSPECIFIED.Descriptor()
	values := enumDesc.Values()
	for i := range values.Len() {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()
		if proto.HasExtension(opts, mbzv1.E_ApiEnumValue) {
			apiValues := proto.GetExtension(opts, mbzv1.E_ApiEnumValue).([]string)
			if slices.Contains(apiValues, typeURI) {
				return mbzv1.ProblemDetail_Type(valueDesc.Number())
			}
		}
	}
	return mbzv1.ProblemDetail_TYPE_UNRECOGNIZED
}

func optString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optInt32(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}

func httpStatusToConnectCode(statusCode int) connect.Code {
	switch statusCode {
	case http.StatusBadRequest:
		return connect.CodeInvalidArgument
	case http.StatusUnauthorized:
		return connect.CodeUnauthenticated
	case http.StatusForbidden:
		return connect.CodePermissionDenied
	case http.StatusNotFound:
		return connect.CodeNotFound
	case http.StatusConflict:
		return connect.CodeAlreadyExists
	case http.StatusTooManyRequests:
		return connect.CodeResourceExhausted
	case http.StatusNotImplemented:
		return connect.CodeUnimplemented
	case http.StatusServiceUnavailable:
		return connect.CodeUnavailable
	case http.StatusGatewayTimeout:
		return connect.CodeDeadlineExceeded
	case http.StatusInternalServerError:
		return connect.CodeInternal
	default:
		return connect.CodeUnknown
	}
}
