package mbz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

// DebugTransport is a [http.RoundTripper] that dumps requests and responses to stderr.
// When Enabled is non-nil, it checks the pointed-to bool on each request,
// allowing lazy binding to a CLI flag.
type DebugTransport struct {
	// Enabled controls whether debug output is active.
	// When nil, debug output is always active.
	Enabled *bool
	// Next is the underlying transport. When nil, [http.DefaultTransport] is used.
	Next http.RoundTripper
}

func (t *DebugTransport) next() http.RoundTripper {
	if t.Next != nil {
		return t.Next
	}
	return http.DefaultTransport
}

func (t *DebugTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.Enabled != nil && !*t.Enabled {
		return t.next().RoundTrip(request)
	}
	requestDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, fmt.Errorf("failed to dump request for debug: %w", err)
	}
	prettyPrintDump(os.Stderr, requestDump, "> ")
	response, err := t.next().RoundTrip(request)
	if err != nil {
		return nil, err
	}
	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		return nil, fmt.Errorf("failed to dump response for debug: %w", err)
	}
	prettyPrintDump(os.Stderr, responseDump, "< ")
	return response, nil
}

func prettyPrintDump(w io.Writer, dump []byte, prefix string) {
	var output bytes.Buffer
	output.Grow(len(dump) * 2)
	for line := range bytes.Lines(dump) {
		output.WriteString(prefix)
		output.Write(line)
	}
	output.WriteByte('\n')
	_, _ = w.Write(output.Bytes())
}
