package trading

// Market represents a trading pair/market
type Market struct {
	Label     string    `json:"label"`
	Symbol    string    `json:"symbol"`
	Base      string    `json:"base"`
	Quote     string    `json:"quote"`
	AssetType AssetType `json:"assetType"` // spot, perpetual, futures, option — see AssetType
	Address   string    `json:"address,omitempty"` // on-chain address (pool, contract, token mint) for DEX sources

	// Precision & limits — all as string to preserve exact value
	PriceTick    string `json:"priceTick"`             // e.g. "0.1", "0.00001"
	QuantityTick string `json:"quantityTick"`          // e.g. "0.00001"
	MinQuantity  string `json:"minQuantity"`           // e.g. "0.00001"
	MaxQuantity  string `json:"maxQuantity,omitempty"` // optional
	MinNotional  string `json:"minNotional,omitempty"` // if price × quantity < minNotional => order rejected (e.g. "1" or "5.0")
	MaxNotional  string `json:"maxNotional,omitempty"` // if price × quantity > maxNotional => order rejected

	// Fees & rates
	MakerFee       string `json:"makerFee,omitempty"` // e.g. "0.0002" (0.02%)
	TakerFee       string `json:"takerFee,omitempty"` // e.g. "0.0005" (0.05%)
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

	// Contract-quantized markets (COIN-M style inverse contracts): quantities
	// are expressed in contracts of ContractSize instead of base asset.
	// QuantityDenomination is "" or QuantityBase for base-asset quantities
	// (the default every current consumer understands); QuantityContracts
	// marks markets whose Quantity/QuantityTick/MinQuantity are contract
	// counts. SettlementAsset names the asset margin and PnL settle in
	// (e.g. "USDT" for USDT-margined linear, the base coin for inverse).
	// Consumers that do not support contract-quantized markets MUST reject
	// them explicitly rather than misinterpret quantities.
	QuantityDenomination QuantityDenomination `json:"quantityDenomination,omitempty"`
	SettlementAsset      string               `json:"settlementAsset,omitempty"`

	// Optional derived UI helpers (safe as int)
	PricePrecision    int `json:"pricePrecision,omitempty"` // derived: -log10(tick)
	QuantityPrecision int `json:"quantityPrecision,omitempty"`

	// Extra carries venue-specific market metadata with no first-class field
	// (e.g. Binance PERCENT_PRICE bands, MAX_NUM_ORDERS) instead of dropping it.
	Extra map[string]any `json:"extra,omitempty"`
}

// Key returns the canonical market identity: assetType + raw venue symbol.
//
// The same raw symbol can exist in several market types on one venue (Binance
// spot BTCUSDT and USDⓈ-M perpetual BTCUSDT), so Symbol alone is NOT a stable
// identity. Consumers that key, cache or cross-reference markets (market
// proxy configs, snapshot grouping, pickers) MUST use this composite key —
// not Symbol.
func (m Market) Key() string {
	return string(m.AssetType) + ":" + m.Symbol
}

// AssetType classifies a market. The vocabulary is CLOSED: every datasource
// plugin uses exactly these values, so consumers can route and group on them.
//
//	"perpetual" is any linear cash-settled perpetual (Binance USDⓈ-M, Bybit
//	linear, Orderly PERP_*), regardless of what the venue itself calls it.
//	"futures" is a DATED delivery future (has ExpiryTimestamp); what venues
//	loosely call "futures" that are actually perpetuals maps to "perpetual".
//
// The inverse_* values are RESERVED for inverse/coin-margined contracts.
// They pair with QuantityDenomination = QuantityContracts and a base-asset
// SettlementAsset. Nothing implements them yet; consumers should treat them
// as unsupported until a plugin does.
type AssetType string

const (
	AssetTypeSpot             AssetType = "spot"
	AssetTypePerpetual        AssetType = "perpetual"
	AssetTypeFutures          AssetType = "futures" // dated delivery (ExpiryTimestamp set)
	AssetTypeOption           AssetType = "option"
	AssetTypeInversePerpetual AssetType = "inverse_perpetual" // reserved
	AssetTypeInverseFutures   AssetType = "inverse_futures"   // reserved
)

// QuantityDenomination declares what a market's quantity fields express.
type QuantityDenomination string

const (
	// QuantityBase is base-asset quantity (the default; empty means the same).
	QuantityBase QuantityDenomination = "base"
	// QuantityContracts is a count of fixed-size contracts (ContractSize);
	// multiply by ContractSize to obtain base-asset amounts.
	QuantityContracts QuantityDenomination = "contracts"
)

// IsDerivative reports whether the market type trades on margin (perpetuals
// and dated futures, linear or inverse). Derivative orders belong in the
// futures response scope; spot orders in the spot scope.
func (a AssetType) IsDerivative() bool {
	switch a {
	case AssetTypePerpetual, AssetTypeFutures, AssetTypeInversePerpetual, AssetTypeInverseFutures:
		return true
	}
	return false
}
