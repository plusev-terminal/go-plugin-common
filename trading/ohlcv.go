package trading

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
