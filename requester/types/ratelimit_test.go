package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func local429Response(retryAfter string) *Response {
	return &Response{
		Status: http.StatusTooManyRequests,
		Headers: http.Header{
			"Retry-After":          []string{retryAfter},
			"X-Terminal-Ratelimit": []string{"local"},
		},
		Error: "rate limited by Terminal: host paused for 7s",
	}
}

func TestErrorFromResponse_Local429(t *testing.T) {
	err := ErrorFromResponse(local429Response("7"))
	if err == nil {
		t.Fatal("Expected an error for a local 429")
	}

	var rle *RateLimitedError
	if !errors.As(err, &rle) {
		t.Fatalf("Expected *RateLimitedError, got %T: %v", err, err)
	}
	if rle.RetryAfter != 7*time.Second {
		t.Fatalf("Expected RetryAfter 7s, got %v", rle.RetryAfter)
	}
	if rle.Error() != "rate limited by Terminal: host paused for 7s" {
		t.Fatalf("Unexpected message: %q", rle.Error())
	}
}

func TestErrorFromResponse_Local429MissingRetryAfter(t *testing.T) {
	res := local429Response("7")
	res.Headers.Del("Retry-After")

	var rle *RateLimitedError
	if !errors.As(ErrorFromResponse(res), &rle) {
		t.Fatal("Expected *RateLimitedError even without Retry-After")
	}
	if rle.RetryAfter != 0 {
		t.Fatalf("Expected RetryAfter 0, got %v", rle.RetryAfter)
	}
}

func TestErrorFromResponse_PlainHostError(t *testing.T) {
	res := &Response{Status: http.StatusInternalServerError, Error: "boom"}

	err := ErrorFromResponse(res)
	if err == nil {
		t.Fatal("Expected an error")
	}
	if err.Error() != "boom" {
		t.Fatalf("Expected message boom, got %q", err.Error())
	}

	var rle *RateLimitedError
	if errors.As(err, &rle) {
		t.Fatal("Plain host error must not be a *RateLimitedError")
	}
}

// The typed error requires BOTH the 429 status and the local marker header.
func TestErrorFromResponse_MarkerRequiresBothSignals(t *testing.T) {
	// Upstream 429 with an error but without the marker header: plain error.
	upstream := &Response{Status: http.StatusTooManyRequests, Error: "upstream busy"}
	err := ErrorFromResponse(upstream)
	var rle *RateLimitedError
	if errors.As(err, &rle) {
		t.Fatal("429 without the local marker must stay a plain error")
	}

	// Marker header on a non-429: plain error.
	wrongStatus := &Response{Status: http.StatusOK, Error: "weird"}
	wrongStatus.Headers = http.Header{"X-Terminal-Ratelimit": []string{"local"}}
	if errors.As(ErrorFromResponse(wrongStatus), &rle) {
		t.Fatal("Marker header without 429 must stay a plain error")
	}
}

func TestErrorFromResponse_NoErrorMeansNil(t *testing.T) {
	if err := ErrorFromResponse(&Response{Status: 200}); err != nil {
		t.Fatalf("Expected nil, got %v", err)
	}

	// Upstream 429 passed through without an Error: normal response.
	upstream := &Response{Status: http.StatusTooManyRequests}
	if err := ErrorFromResponse(upstream); err != nil {
		t.Fatalf("Expected nil for an upstream 429 without Error, got %v", err)
	}
}

func TestIsRateLimited(t *testing.T) {
	err := ErrorFromResponse(local429Response("7"))
	retryAfter, ok := IsRateLimited(fmt.Errorf("send failed: %w", err))
	if !ok || retryAfter != 7*time.Second {
		t.Fatalf("Expected (7s, true), got (%v, %v)", retryAfter, ok)
	}

	retryAfter, ok = IsRateLimited(errors.New("boom"))
	if ok || retryAfter != 0 {
		t.Fatalf("Expected (0, false), got (%v, %v)", retryAfter, ok)
	}

	retryAfter, ok = IsRateLimited(nil)
	if ok || retryAfter != 0 {
		t.Fatalf("Expected (0, false) for nil, got (%v, %v)", retryAfter, ok)
	}
}

func TestRequestJSON_WithoutHints(t *testing.T) {
	// All four v1.14.0 fields are always emitted; the hints must not appear.
	data, err := json.Marshal(&Request{URL: "https://x/", Method: "GET"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	want := `{"url":"https://x/","method":"GET","headers":null,"body":null}`
	if string(data) != want {
		t.Fatalf("Expected %s, got %s", want, string(data))
	}
}

func TestRequestJSON_WithHints(t *testing.T) {
	data, err := json.Marshal(&Request{URL: "https://rpc/x", Method: "POST", RateKey: "eth_call", Cost: 120})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	want := `{"url":"https://rpc/x","method":"POST","headers":null,"body":null,"rate_key":"eth_call","cost":120}`
	if string(data) != want {
		t.Fatalf("Expected %s, got %s", want, string(data))
	}
}
