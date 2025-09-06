package datasrc

// CredentialField represents a single credential field required for authentication when interacting with the 3rd party api.
type CredentialField struct {
	Name      string `json:"name"`
	Encrypt   bool   `json:"encrypt"`
	Mask      bool   `json:"mask"`
	OmitEmpty bool   `json:"omitEmpty"`
}

// MarketMeta represents metadata about a trading market/pair
type MarketMeta struct {
	Name      string `json:"name"`      // Trading pair name (e.g., "BTCUSDT")
	Base      string `json:"base"`      // Base asset (e.g., "BTC")
	Quote     string `json:"quote"`     // Quote asset (e.g., "USDT")
	AssetType string `json:"assetType"` // Asset type (e.g., "spot", "futures")
}

// Timeframe represents a supported timeframe for OHLCV data
type Timeframe struct {
	Label    string `json:"label"`    // Human-readable label (e.g., "1m", "5m")
	ApiValue string `json:"apiValue"` // Value used for API calls
	Interval int64  `json:"interval"` // Interval in seconds
}

// OHLCVParams represents parameters for OHLCV data requests
type OHLCVParams struct {
	Credentials map[string]string `json:"credentials"`
	Symbol      string            `json:"symbol"`    // Trading pair symbol
	Timeframe   string            `json:"timeframe"` // Timeframe for the data
	StartTime   int64             `json:"startTime"` // Start timestamp (Unix)
	EndTime     int64             `json:"endTime"`   // End timestamp (Unix)
	Limit       int               `json:"limit"`     // Maximum number of records
}

// OHLCVRecord represents a single OHLCV (candlestick) data point
type OHLCVRecord struct {
	Timestamp int64   `json:"timestamp"` // Unix timestamp
	Open      float64 `json:"open"`      // Opening price
	High      float64 `json:"high"`      // Highest price
	Low       float64 `json:"low"`       // Lowest price
	Close     float64 `json:"close"`     // Closing price
	Volume    float64 `json:"volume"`    // Trading volume
}

// StreamSetupRequest represents the request sent to plugin for stream setup
type StreamSetupRequest struct {
	StreamID   string                 `json:"streamId"`
	StreamType string                 `json:"streamType"` // "ohlcv", "orderbook", "orders", "trades", etc.
	Parameters map[string]interface{} `json:"parameters"` // Generic parameters
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
	Success     bool        `json:"success"`
	Action      string      `json:"action"`             // "ignore", "data", "reconnect", "close", "send"
	DataType    string      `json:"dataType,omitempty"` // "ohlcv", "orderbook", "order_fill", etc.
	Data        interface{} `json:"data,omitempty"`     // Generic data payload
	SendMessage string      `json:"sendMessage,omitempty"`
	Error       string      `json:"error,omitempty"`
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
