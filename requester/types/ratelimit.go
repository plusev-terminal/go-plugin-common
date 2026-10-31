package types

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

// rateLimitHeader marks a 429 as produced by the host's local rate limiter
// (as opposed to an upstream 429 passed through to the plugin).
const rateLimitHeader = "X-Terminal-Ratelimit"

// RateLimitedError is returned by Send when the host's local rate limiter
// rejected the call before any bytes hit the wire. RetryAfter is how long
// the host asks the caller to wait before retrying; Message is the host's
// error text.
type RateLimitedError struct {
	RetryAfter time.Duration
	Message    string
}

// Error returns the host's error message.
func (e *RateLimitedError) Error() string {
	return e.Message
}

// IsRateLimited reports whether err was produced by the host's local rate
// limiter, and if so how long to wait before retrying.
func IsRateLimited(err error) (retryAfter time.Duration, ok bool) {
	var rle *RateLimitedError
	if !errors.As(err, &rle) {
		return 0, false
	}
	return rle.RetryAfter, true
}

// ErrorFromResponse maps a host response onto an error, or nil when the
// plugin should treat the response as a normal result. The host's local
// rate limiter answers with Status 429, the X-Terminal-Ratelimit: local
// header, a Retry-After header in whole seconds and an Error message; that
// exact combination becomes a *RateLimitedError. Any other non-empty Error
// becomes a plain error (the pre-existing behavior); responses without an
// Error map to nil — including upstream 429s, which the plugin handles
// itself from the returned response.
func ErrorFromResponse(res *Response) error {
	if res.Error == "" {
		return nil
	}

	if res.Status == http.StatusTooManyRequests && res.Headers.Get(rateLimitHeader) == "local" {
		return &RateLimitedError{
			RetryAfter: retryAfterFromHeader(res.Headers),
			Message:    res.Error,
		}
	}

	return errors.New(res.Error)
}

// retryAfterFromHeader parses a Retry-After delta-seconds header; a
// missing or unparseable header yields 0.
func retryAfterFromHeader(h http.Header) time.Duration {
	secs, err := strconv.Atoi(h.Get("Retry-After"))
	if err != nil || secs < 0 {
		return 0
	}
	return time.Duration(secs) * time.Second
}
