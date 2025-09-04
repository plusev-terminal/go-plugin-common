package datasrc

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
	Symbol    string `json:"symbol"`    // Trading pair symbol
	Timeframe string `json:"timeframe"` // Timeframe for the data
	StartTime int64  `json:"startTime"` // Start timestamp (Unix)
	EndTime   int64  `json:"endTime"`   // End timestamp (Unix)
	Limit     int    `json:"limit"`     // Maximum number of records
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

// StreamConfig represents configuration for streaming data
type StreamConfig struct {
	Symbol   string `json:"symbol"`   // Trading pair symbol
	Interval int64  `json:"interval"` // Update interval in seconds
}
