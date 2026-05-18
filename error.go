package mbz

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"connectrpc.com/connect"
)

func newResponseError(httpResponse *http.Response) error {
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		body = fmt.Appendf(nil, "failed to read response body: %s", err)
	}
	var msg string
	if len(body) > 0 {
		msg = fmt.Sprintf("http %d: %s", httpResponse.StatusCode, body)
	} else {
		msg = fmt.Sprintf("http %d", httpResponse.StatusCode)
	}
	return connect.NewError(httpStatusToConnectCode(httpResponse.StatusCode), errors.New(msg))
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
