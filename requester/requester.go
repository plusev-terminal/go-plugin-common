package requester

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/extism/go-pdk"
)

//go:wasmimport extism:host/user http_request
func httpRequest(uint64) uint64

// Request is the request to be sent to the host
type Request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

// Response is the response from the host
type Response struct {
	StatusCode int         `json:"statusCode"`
	Headers    http.Header `json:"headers"`
	Body       []byte      `json:"body"`
	Error      string      `json:"error,omitempty"`
}

// Send sends the request to the host and returns the response.
// If v is not nil, the response body will be unmarshaled into it.
func Send(req *Request, v any) (*Response, error) {
	mem, err := pdk.AllocateJSON(req)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for request: %w", err)
	}

	ptr := httpRequest(mem.Offset())
	rmem := pdk.FindMemory(ptr)
	respData := rmem.ReadBytes()

	var res Response
	if err := json.Unmarshal(respData, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if res.Error != "" {
		return nil, errors.New(res.Error)
	}

	if v != nil {
		if err := json.Unmarshal(res.Body, v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body into target struct: %w", err)
		}
	}

	return &res, nil
}
