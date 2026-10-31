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
	// RateKey is the bucket value for PerKey rate limit rules, e.g. a
	// JSON-RPC method or endpoint group. When empty, the host derives it
	// from the JSON-RPC "method" field of the body itself — but only if a
	// PerKey rule applies to the request.
	RateKey string `json:"rate_key,omitempty"`
	// Cost is the number of units the request consumes for Weighted rate
	// limit rules, e.g. an exchange's request weight. 0 counts as 1.
	Cost int `json:"cost,omitempty"`
}

// Response is the response from the host
type Response struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
	Error   string      `json:"error,omitempty"`
}
