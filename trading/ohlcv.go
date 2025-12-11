package trading

// OHLCVRecord represents a single OHLCV (candlestick) data point
// Price and volume fields are strings to preserve precision for tokens
// with very small values (e.g., 0.000000123456). Consumers should use
// high-precision libraries like shopspring/decimal to parse these values.
type OHLCVRecord struct {
	OpenTime int64  `json:"openTime"`
	Open     string `json:"open"`
	High     string `json:"high"`
	Low      string `json:"low"`
	Close    string `json:"close"`
	Volume   string `json:"volume"`
}
