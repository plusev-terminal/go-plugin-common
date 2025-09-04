package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	rt "github.com/plusev-terminal/go-plugin-common/requester/types"
)

// MockRequester implements requester.Interface for testing using standard net/http
// This allows testing plugins without the need for WASM host functions
type MockRequester struct {
	client    *http.Client
	responses map[string]string // URL pattern -> JSON response
	errors    map[string]error  // URL pattern -> error
	calls     []string          // Track all calls made
}

// NewMockRequester creates a new mock requester for testing
func NewMockRequester() *MockRequester {
	return &MockRequester{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		responses: make(map[string]string),
		errors:    make(map[string]error),
		calls:     make([]string, 0),
	}
}

// SetMockResponse sets a mock JSON response for a URL pattern
// Use patterns like "/v3/public/instruments" or wildcards like "/v3/public/*"
func (m *MockRequester) SetMockResponse(urlPattern string, jsonResponse string) {
	m.responses[urlPattern] = jsonResponse
}

// SetMockError sets a mock error for a URL pattern
func (m *MockRequester) SetMockError(urlPattern string, err error) {
	m.errors[urlPattern] = err
}

// Send implements requester.Interface for testing
func (m *MockRequester) Send(req *rt.Request, response interface{}) (*rt.Response, error) {
	m.calls = append(m.calls, req.URL)

	// Check for mock errors first
	for pattern, err := range m.errors {
		if matchesPattern(req.URL, pattern) {
			return nil, err
		}
	}

	// Check for mock responses
	for pattern, jsonResp := range m.responses {
		if matchesPattern(req.URL, pattern) {
			// Unmarshal the JSON response into the provided response interface
			if response != nil {
				if err := json.Unmarshal([]byte(jsonResp), response); err != nil {
					return nil, fmt.Errorf("failed to unmarshal mock response: %w", err)
				}
			}

			return &rt.Response{
				Status:  200,
				Headers: http.Header{"Content-Type": []string{"application/json"}},
				Body:    []byte(jsonResp),
			}, nil
		}
	}

	// If no mock is set, make a real HTTP request (useful for integration tests)
	return m.makeRealRequest(req, response)
}

// makeRealRequest makes an actual HTTP request using net/http
// This is useful for integration testing against real APIs
func (m *MockRequester) makeRealRequest(req *rt.Request, response interface{}) (*rt.Response, error) {
	var body io.Reader
	if len(req.Body) > 0 {
		body = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := m.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// If response interface is provided, unmarshal the JSON
	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return &rt.Response{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Body:    respBody,
	}, nil
}

// GetCalls returns all URLs that were called during testing
func (m *MockRequester) GetCalls() []string {
	return m.calls
}

// Reset clears all mock responses and call history
func (m *MockRequester) Reset() {
	m.responses = make(map[string]string)
	m.errors = make(map[string]error)
	m.calls = make([]string, 0)
}

// matchesPattern checks if URL matches the pattern
func matchesPattern(url, pattern string) bool {
	if url == pattern {
		return true
	}
	// Support wildcard patterns ending with *
	if strings.HasSuffix(pattern, "*") {
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(url, prefix)
	}
	// Check if pattern is contained in URL
	return strings.Contains(url, pattern)
}
