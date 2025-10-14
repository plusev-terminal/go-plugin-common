package datasrc

import "time"

// Command represents a request to a data source
type Command struct {
	Name   string         `json:"name"`   // e.g., "ohlcvStream", "getMarkets", "getBalance"
	Params map[string]any `json:"params"` // Flexible parameters specific to each command
}

// Response represents the result of a command execution
type Response struct {
	Result          bool   `json:"result"`
	Data            any    `json:"data,omitempty"`            // Could be direct data or a channel for streams
	Error           string `json:"error,omitempty"`           // Error message if Success is false
	CacheForSeconds *int64 `json:"cacheForSeconds,omitempty"` // Optional: cache duration in seconds (wrapper converts to time.Duration)
}

// StreamData represents a single piece of data from a stream
type StreamData struct {
	StreamID string `json:"streamId"` // Unique identifier for this stream
	Data     any    `json:"data"`     // The actual data (e.g., OHLCV candle, orderbook update)
}

// ConfigField defines a configuration field that a data source requires
// This is used to generate UI forms for setting up connections
type ConfigField struct {
	Name        string `json:"name"`                  // Field name (e.g., "apiKey", "applicationID")
	Label       string `json:"label"`                 // Human-readable label for UI
	Type        string `json:"type"`                  // Input type: "text", "password", "number", etc.
	Required    bool   `json:"required"`              // Whether this field is mandatory
	Encrypt     bool   `json:"encrypt"`               // Whether to encrypt this field in database
	Mask        bool   `json:"mask"`                  // Whether to mask this field in API responses
	Placeholder string `json:"placeholder,omitempty"` // Placeholder text for UI
	Description string `json:"description,omitempty"` // Help text explaining the field
}

// Market represents a trading pair/market
type Market struct {
	Symbol    string `json:"symbol"`    // Exchange-specific symbol (e.g., "BTCUSDT", "BTC/USDT")
	Base      string `json:"base"`      // Base currency (e.g., "BTC")
	Quote     string `json:"quote"`     // Quote currency (e.g., "USDT")
	AssetType string `json:"assetType"` // e.g., "spot", "futures", "perpetual"
}

// OHLCVParams represents parameters for OHLCV data requests
type OHLCVParams struct {
	Symbol    string `json:"symbol"`    // Trading pair symbol
	Timeframe string `json:"timeframe"` // Timeframe for the data
	StartTime int64  `json:"startTime"` // Start timestamp (Unix)
	EndTime   int64  `json:"endTime"`   // End timestamp (Unix)
	Limit     int    `json:"limit"`     // Maximum number of records
}

// OHLCVRecord represents a single OHLCV (candlestick) data point
// Price and volume fields are strings to preserve precision for tokens
// with very small values (e.g., 0.000000123456). Consumers should use
// high-precision libraries like shopspring/decimal to parse these values.
type OHLCVRecord struct {
	Timestamp int64  `json:"timestamp"` // Unix timestamp
	Open      string `json:"open"`      // Opening price (string for arbitrary precision)
	High      string `json:"high"`      // Highest price (string for arbitrary precision)
	Low       string `json:"low"`       // Lowest price (string for arbitrary precision)
	Close     string `json:"close"`     // Closing price (string for arbitrary precision)
	Volume    string `json:"volume"`    // Trading volume (string for arbitrary precision)
}

// StreamSetupRequest represents the request sent to plugin for stream setup
type StreamSetupRequest struct {
	StreamID   string         `json:"streamId"`
	StreamType string         `json:"streamType"` // "ohlcv", "orderbook", "orders", "trades", etc.
	Parameters map[string]any `json:"parameters"` // Generic parameters
}

// StreamSetupResponse represents plugin's response to stream setup request
type StreamSetupResponse struct {
	Success         bool              `json:"success"`
	WebSocketURL    string            `json:"websocketUrl"`
	Headers         map[string]string `json:"headers,omitempty"`
	Subprotocol     string            `json:"subprotocol,omitempty"`
	InitialMessages []string          `json:"initialMessages"`
	Error           string            `json:"error,omitempty"`
}

// StreamMessageRequest represents the request sent to plugin for message processing
type StreamMessageRequest struct {
	StreamID     string `json:"streamId"`
	ConnectionID string `json:"connectionId"`
	Message      string `json:"message"`
	MessageType  string `json:"messageType"` // "data", "error", "close"
}

// StreamMessageResponse represents plugin's response to a stream message
type StreamMessageResponse struct {
	Success     bool   `json:"success"`
	Action      string `json:"action"`             // "ignore", "data", "reconnect", "close", "send"
	DataType    string `json:"dataType,omitempty"` // "ohlcv", "orderbook", "order_fill", etc.
	Data        any    `json:"data,omitempty"`     // Generic data payload
	SendMessage string `json:"sendMessage,omitempty"`
	Error       string `json:"error,omitempty"`
}

// StreamConnectionEvent represents a connection lifecycle event
type StreamConnectionEvent struct {
	StreamID     string `json:"streamId"`
	ConnectionID string `json:"connectionId"`
	EventType    string `json:"eventType"` // "connected", "disconnected", "error"
	Error        string `json:"error,omitempty"`
}

// StreamConnectionResponse represents plugin's response to a connection event
type StreamConnectionResponse struct {
	Success bool   `json:"success"`
	Action  string `json:"action"` // "ignore", "reconnect", "close"
	Error   string `json:"error,omitempty"`
}

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
	Command string         `json:"command"` // Command name or "*" for wildcard
	Scope   RateLimitScope `json:"scope"`   // Use RateLimitScope constants
	RPS     float64        `json:"rps"`     // Requests per second (can be fractional, e.g., 0.1 = 1 req per 10 sec)
	Burst   int            `json:"burst"`   // Burst allowance
	Cost    int            `json:"cost"`    // Token cost per request (default: 1, for commands that make multiple API calls)
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
