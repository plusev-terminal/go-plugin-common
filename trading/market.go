package trading

// Market represents a trading pair/market
type Market struct {
	Label     string `json:"label"`
	Symbol    string `json:"symbol"`
	Base      string `json:"base"`
	Quote     string `json:"quote"`
	AssetType string `json:"assetType"` // "spot", "perpetual", "futures", "option"

	// Precision & limits — all as string to preserve exact value
	PriceTick    string `json:"priceTick"`             // e.g. "0.1", "0.00001"
	QuantityTick string `json:"quantityTick"`          // e.g. "0.00001"
	MinQuantity  string `json:"minQuantity"`           // e.g. "0.00001"
	MaxQuantity  string `json:"maxQuantity,omitempty"` // optional
	MinNotional  string `json:"minNotional,omitempty"` // if price × quantity < minNotional => order rejected (e.g. "1" or "5.0")
	MaxNotional  string `json:"maxNotional,omitempty"` // if price × quantity > maxNotional => order rejected

	// Fees & rates
	MakerFee       string `json:"makerFee,omitempty"` // e.g. "0.0002" (0.02%)
	TakerFee       string `json:"takerFee,omitempty"` // e.g. "0.0005"
	LiquidationFee string `json:"liquidationFee,omitempty"`

	// Leverage & margin
	MaxLeverage           string `json:"maxLeverage,omitempty"`       // e.g. "100" or "125"
	InitialMarginRate     string `json:"initialMarginRate,omitempty"` // e.g. "0.01" → 100x
	MaintenanceMarginRate string `json:"maintenanceMarginRate,omitempty"`

	// Funding (perpetuals)
	FundingInterval int    `json:"fundingInterval,omitempty"` // e.g. 8 (hours)
	FundingCap      string `json:"fundingCap,omitempty"`      // e.g. "0.000375"
	FundingFloor    string `json:"fundingFloor,omitempty"`

	// Other common fields
	ContractSize    string `json:"contractSize,omitempty"`    // e.g. "1" for linear, "0.0001" for inverse
	ExpiryTimestamp int64  `json:"expiryTimestamp,omitempty"` // 0 for perps
	Status          string `json:"status,omitempty"`          // "TRADING", "HALTED", etc.

	// Optional derived UI helpers (safe as int)
	PricePrecision    int `json:"pricePrecision,omitempty"` // derived: -log10(tick)
	QuantityPrecision int `json:"quantityPrecision,omitempty"`
}
