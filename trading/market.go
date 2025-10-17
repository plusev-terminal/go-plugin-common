package trading

// Market represents a trading pair/market
type Market struct {
	Label     string `json:"label"`     // Human-readable label (e.g., "BTC/USDT")
	Symbol    string `json:"symbol"`    // Exchange-specific symbol (e.g., "BTCUSDT", "BTC/USDT", "PERP_BTC_USDT")
	Base      string `json:"base"`      // Base currency (e.g., "BTC")
	Quote     string `json:"quote"`     // Quote currency (e.g., "USDT")
	AssetType string `json:"assetType"` // e.g., "spot", "futures", "perpetual"
}
