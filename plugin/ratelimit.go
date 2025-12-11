package plugin

import "time"

// RateLimitScope defines the scope at which rate limiting is enforced
type RateLimitScope string

const (
	// RateLimitScopeIP applies rate limiting per IP address (used for public endpoints and IP-based limits)
	RateLimitScopeIP RateLimitScope = "ip"

	// RateLimitScopeAPIKey applies rate limiting per API key (used for authenticated endpoints)
	RateLimitScopeAPIKey RateLimitScope = "apikey"
)

// RateLimit defines rate limit configuration for a command
type RateLimit struct {
	Command string           `json:"command"` // Command name or "*" for wildcard
	Scope   []RateLimitScope `json:"scope"`   // Scope keys (e.g., []RateLimitScope{RateLimitScopeIP}, []RateLimitScope{RateLimitScopeAPIKey, "user_id"})
	RPS     float64          `json:"rps"`     // Requests per second (can be fractional, e.g., 0.1 = 1 req per 10 sec)
	Burst   int              `json:"burst"`   // Burst allowance
	Cost    int              `json:"cost"`    // Token cost per request (default: 1, for commands that make multiple API calls)
}

// CalculateRPS converts a request count and time duration to requests per second.
// This helper makes it easy to define rate limits like "100 requests per minute"
// without doing the math yourself.
//
// Example usage:
//
//	CalculateRPS(100, time.Minute)        // 100 requests per minute = 1.666... RPS
//	CalculateRPS(6000, time.Minute)       // 6000 requests per minute = 100 RPS
//	CalculateRPS(1, 10*time.Second)       // 1 request per 10 seconds = 0.1 RPS
//	CalculateRPS(61000, 5*time.Minute)    // 61000 requests per 5 minutes = 203.333... RPS
func CalculateRPS(requests int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	seconds := duration.Seconds()
	if seconds == 0 {
		return 0
	}

	return float64(requests) / seconds
}
