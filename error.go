package mbz

import (
	"fmt"
	"io"
	"net/http"

	"connectrpc.com/connect"
)

// ResponseError represents a Mercedes-Benz API response error.
type ResponseError struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Body is the body of the response.
	Body []byte
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
	return respErr.connectError()
}

// Error implements the error interface.
func (e *ResponseError) Error() string {
	if len(e.Body) > 0 {
		return fmt.Sprintf("http %d: %s", e.StatusCode, string(e.Body))
	}
	return fmt.Sprintf("http %d", e.StatusCode)
}

func (e *ResponseError) connectError() error {
	return connect.NewError(httpStatusToConnectCode(e.StatusCode), e)
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
