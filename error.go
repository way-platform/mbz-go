package mbz

import (
	"fmt"
	"net/http"
)

// newResponseError creates a new error from an HTTP response.
func newResponseError(httpResponse *http.Response) error {
	// TOFDO: Implement error parsing.
	return fmt.Errorf("HTTP %s", httpResponse.Status)
}
