package types

import "net/http"

// RequestDoer defines the contract for making HTTP requests
// This allows for easy mocking and testing of plugins
type RequestDoer interface {
	Send(req *Request, v any) (*Response, error)
}

// Request is the request to be sent to the host
type Request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

// Response is the response from the host
type Response struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
	Error   string      `json:"error,omitempty"`
}
