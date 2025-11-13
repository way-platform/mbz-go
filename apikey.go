package mbz

import "net/http"

// apiKeyTransport is a [http.RoundTripper] that sets the x-api-key header on the request.
type apiKeyTransport struct {
	next   http.RoundTripper
	apiKey string
}

// RoundTrip implements the [http.RoundTripper] interface.
func (t *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("x-api-key", t.apiKey)
	return t.next.RoundTrip(req)
}
