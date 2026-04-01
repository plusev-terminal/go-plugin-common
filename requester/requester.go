package requester

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/extism/go-pdk"
	rt "github.com/plusev-terminal/go-plugin-common/requester/types"
)

//go:wasmimport extism:host/user http_request
func httpRequest(uint64) uint64

// Requester is the default requester that uses the host functions
type Requester struct{}

// NewRequester creates a new default requester
func NewRequester() *Requester {
	return &Requester{}
}

// Send sends the request to the host and returns the response.
// If v is not nil, the response body will be unmarshaled into it.
func (d *Requester) Send(req *rt.Request, v any) (*rt.Response, error) {
	mem, err := pdk.AllocateJSON(req)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for request: %w", err)
	}

	ptr := httpRequest(mem.Offset())
	rmem := pdk.FindMemory(ptr)
	respData := rmem.ReadBytes()

	var res rt.Response
	if err := json.Unmarshal(respData, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if res.Error != "" {
		return nil, errors.New(res.Error)
	}

	if v != nil {
		// Guard against non-JSON responses (e.g. HTML error/maintenance pages)
		// which can cause excessive memory allocation during json.Unmarshal
		// and crash the WASM module with "out of bounds memory access".
		ct := res.Headers.Get("Content-Type")
		if ct != "" && !strings.Contains(ct, "json") {
			return &res, fmt.Errorf("unexpected Content-Type %q (expected JSON), status %d, body prefix: %.120s", ct, res.Status, string(res.Body))
		}

		if err := json.Unmarshal(res.Body, v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body into target struct: %w", err)
		}
	}

	return &res, nil
}
