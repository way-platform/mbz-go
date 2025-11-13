package mbz

import (
	"fmt"
	"io"
	"net/http"
)

// ResponseError represents a VW Fleet Interface response error.
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
	return &ResponseError{
		StatusCode: httpResponse.StatusCode,
		Body:       body,
	}
}

// Error implements the error interface.
func (e *ResponseError) Error() string {
	if len(e.Body) > 0 {
		return fmt.Sprintf("http %d: %s", e.StatusCode, string(e.Body))
	}
	return fmt.Sprintf("http %d", e.StatusCode)
}
