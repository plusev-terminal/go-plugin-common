package trading

import (
	"fmt"
	"strings"
)

// Market represents a trading pair/market
type Market struct {
	Label     string    `json:"label"`
	Symbol    string    `json:"symbol"`
	Base      string    `json:"base"`
	Quote     string    `json:"quote"`
	AssetType AssetType `json:"assetType"`         // spot, perpetual, futures, option — see AssetType
	Address   string    `json:"address,omitempty"` // on-chain address (pool, contract, token mint) for DEX sources

	// Precision & limits — all as string to preserve exact value. Quantities
	// are in the unit QuantityDenomination declares (base asset by default).
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

	// ContractSize is the size of ONE contract, expressed in ContractSizeAsset
	// (the base asset when ContractSizeAsset is empty). It defines what a
	// quantity of 1 means when QuantityDenomination is QuantityContracts —
	// e.g. OKX BTC-USDT-SWAP "0.01" (BTC, linear) or Binance COIN-M
	// BTCUSD_PERP "100" with ContractSizeAsset "USD" (inverse).
	ContractSize      string `json:"contractSize,omitempty"`
	ContractSizeAsset string `json:"contractSizeAsset,omitempty"`

	ExpiryTimestamp int64  `json:"expiryTimestamp,omitempty"` // 0 for perps
	Status          string `json:"status,omitempty"`          // "TRADING", "HALTED", etc.

	// QuantityDenomination declares what the market's quantity fields — and
	// every order and position quantity on this market — count: base-asset
	// units ("" or QuantityBase, the default) or contracts of ContractSize
	// (QuantityContracts). Contract sizing is NOT an inverse-only concept:
	// OKX, KuCoin and Gate size their linear perpetuals in contracts too.
	// Consumers that size orders and do not handle contract quantities MUST
	// reject such markets (see UsesContracts) instead of misreading them.
	QuantityDenomination QuantityDenomination `json:"quantityDenomination,omitempty"`
	// SettlementAsset names the asset margin and PnL settle in ("USDT" for
	// USDT-margined linear contracts, the base coin for inverse ones).
	SettlementAsset string `json:"settlementAsset,omitempty"`

	// Optional derived UI helpers (safe as int)
	PricePrecision    int `json:"pricePrecision,omitempty"` // derived: -log10(tick)
	QuantityPrecision int `json:"quantityPrecision,omitempty"`
}

// Key returns the canonical identity of a market WITHIN one venue/connection:
// "<assetType>:<symbol>", plus "@<address>" when the market has an on-chain
// Address (DEX venues can list several pools under one symbol).
//
// The same raw symbol can exist in several market types on one venue (Binance
// spot BTCUSDT and USDⓈ-M perpetual BTCUSDT), so Symbol alone is NOT a stable
// identity. Consumers that key, cache or cross-reference markets (proxy
// names, candle caches, snapshot grouping, pickers) MUST use this key — and
// qualify it with the connection when markets of several venues meet.
func (m Market) Key() string {
	key := string(m.AssetType) + ":" + m.Symbol

	if m.Address != "" {
		key += "@" + m.Address
	}

	return key
}

// UsesContracts reports whether quantities on this market count contracts of
// ContractSize instead of base-asset units.
func (m Market) UsesContracts() bool {
	return m.QuantityDenomination == QuantityContracts
}

// FindMarket resolves a market reference — typically a Market a consumer
// persisted (symbol + assetType [+ address]) — against a venue's getMarkets
// list. Symbol (or, for DEX pools, Address) must match; AssetType and Address
// narrow the match when set.
//
// It errors when nothing matches and when the reference is ambiguous (symbol
// only, while the venue lists BTCUSDT as spot AND perpetual): a consumer must
// never guess between market types — they trade at different prices.
func FindMarket(markets []Market, ref Market) (Market, error) {
	if ref.Symbol == "" && ref.Address == "" {
		return Market{}, fmt.Errorf("market reference needs a symbol or an address")
	}

	var found []Market

	for _, m := range markets {
		if ref.Symbol != "" && m.Symbol != ref.Symbol {
			continue
		}

		if ref.AssetType != "" && m.AssetType != ref.AssetType {
			continue
		}

		if ref.Address != "" && m.Address != ref.Address {
			continue
		}

		found = append(found, m)
	}

	switch len(found) {
	case 0:
		return Market{}, fmt.Errorf("market %s not found", describeMarketRef(ref))
	case 1:
		return found[0], nil
	}

	keys := make([]string, 0, len(found))

	for _, m := range found {
		keys = append(keys, m.Key())
	}

	return Market{}, fmt.Errorf("market %q is ambiguous (%s): specify assetType", ref.Symbol, strings.Join(keys, ", "))
}

// describeMarketRef renders a (possibly partial) market reference for errors.
func describeMarketRef(ref Market) string {
	desc := ref.Symbol

	if ref.AssetType != "" {
		desc = string(ref.AssetType) + " " + desc
	}

	if ref.Address != "" {
		desc = strings.TrimSpace(desc + " @" + ref.Address)
	}

	return strings.TrimSpace(desc)
}

// AssetType classifies a market. The vocabulary is CLOSED: every datasource
// plugin uses exactly these values (Valid reports membership), so consumers
// can route and group on them.
//
//	"perpetual" is any linear cash-settled perpetual (Binance USDⓈ-M, Bybit
//	linear, Orderly PERP_*), regardless of what the venue itself calls it.
//	"futures" is a DATED delivery future (has ExpiryTimestamp); what venues
//	loosely call "futures" that are actually perpetuals maps to "perpetual".
//
// The inverse_* values are RESERVED for inverse/coin-margined contracts
// (SettlementAsset = the base coin, ContractSize usually in the quote asset).
// Nothing implements them yet; consumers should treat them as unsupported
// until a plugin does.
//
// The empty AssetType means "unspecified". It is not Valid: only a venue that
// lists a single market type can resolve it (see ResolveAssetType), and order
// requests always require a valid AssetType.
type AssetType string

const (
	AssetTypeSpot             AssetType = "spot"
	AssetTypePerpetual        AssetType = "perpetual"
	AssetTypeFutures          AssetType = "futures" // dated delivery (ExpiryTimestamp set)
	AssetTypeOption           AssetType = "option"
	AssetTypeInversePerpetual AssetType = "inverse_perpetual" // reserved
	AssetTypeInverseFutures   AssetType = "inverse_futures"   // reserved
)

// AssetTypes returns the complete contract vocabulary in a stable order.
func AssetTypes() []AssetType {
	return []AssetType{
		AssetTypeSpot,
		AssetTypePerpetual,
		AssetTypeFutures,
		AssetTypeOption,
		AssetTypeInversePerpetual,
		AssetTypeInverseFutures,
	}
}

// Valid reports whether a is part of the contract vocabulary (the reserved
// inverse_* values included). The empty AssetType is not valid.
func (a AssetType) Valid() bool {
	for _, known := range AssetTypes() {
		if a == known {
			return true
		}
	}

	return false
}

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

// ResolveAssetType is the routing rule every venue applies to market-data and
// order requests: it returns the market type a request for `requested` must
// be served from, given the market types the venue lists.
//
//   - a valid type the venue lists → that type
//   - empty, and the venue lists exactly one type → that type (unambiguous)
//   - empty, and the venue lists several types → error: the caller must say
//     which; a venue never silently defaults (spot and perpetual BTCUSDT are
//     different markets with different prices)
//   - unknown types, or types the venue does not list → error
//
// With no listed types every valid type passes and empty is an error.
func ResolveAssetType(requested AssetType, listed ...AssetType) (AssetType, error) {
	if requested == "" {
		if len(listed) == 1 {
			return listed[0], nil
		}

		if len(listed) == 0 {
			return "", fmt.Errorf("assetType is required")
		}

		return "", fmt.Errorf("assetType is required: this venue lists %s markets", joinAssetTypes(listed))
	}

	if !requested.Valid() {
		return "", fmt.Errorf("unknown assetType %q (want one of %s)", requested, joinAssetTypes(AssetTypes()))
	}

	if len(listed) == 0 {
		return requested, nil
	}

	for _, l := range listed {
		if l == requested {
			return requested, nil
		}
	}

	return "", fmt.Errorf("this venue has no %s markets (it lists %s)", requested, joinAssetTypes(listed))
}

func joinAssetTypes(types []AssetType) string {
	parts := make([]string, 0, len(types))

	for _, t := range types {
		parts = append(parts, string(t))
	}

	return strings.Join(parts, ", ")
}

// QuantityDenomination declares what a market's quantity fields express.
type QuantityDenomination string

const (
	// QuantityBase is base-asset quantity (the default; empty means the same).
	QuantityBase QuantityDenomination = "base"
	// QuantityContracts is a count of fixed-size contracts: one contract is
	// ContractSize units of ContractSizeAsset (base asset when empty).
	QuantityContracts QuantityDenomination = "contracts"
)
