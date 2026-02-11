package account

import (
	"encoding/json"
	"time"
)

// GetPositionsRequest are the optional params for the getPositions command.
// Implementations may use Scopes to avoid unnecessary API calls.
//
// Scopes should refer to balance/portfolio partitions (spot/margin/futures/cfd/etc).
// If omitted, the plugin should return what it can by default.
type GetPositionsRequest struct {
	ScopeFilter
}

// PositionsResponse is the response type for the getPositions command.
// It intentionally contains only position-related data, plus metadata that applies
// to the positions snapshot (e.g., FetchedAt).
//
// Positions are grouped by scope so consumers can display/aggregate without having to
// infer scope from symbol naming.
type PositionsResponse struct {
	FetchedAt time.Time `json:"fetchedAt,omitempty"`

	Scopes map[ScopeType]PositionScope `json:"scopes,omitempty"`

	// Raw stores the unmodified exchange payload for debugging/support.
	Raw json.RawMessage `json:"raw,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// PositionScope represents positions within a specific scope.
type PositionScope struct {
	ScopeID string `json:"scopeId,omitempty"`

	// Network identifies the chain/network when ScopeType is "wallet" (self-custody).
	// Examples: "ethereum", "arbitrum", "solana", "bitcoin".
	Network Network `json:"network,omitempty"`

	// ChainID is the numeric EVM chain id when applicable (e.g., 1 for Ethereum mainnet).
	// Leave as 0 when not applicable or unknown.
	ChainID int `json:"chainId,omitempty"`

	Positions []Position `json:"positions,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// Position represents a normalized open position on an exchange/broker.
// All numeric values are strings to preserve precision.
//
// Quantity should be absolute; use Side to determine direction.
// Side should be one of: "long", "short" or exchange-specific.
type Position struct {
	Symbol string `json:"symbol,omitempty"` // e.g. "BTCUSDT", "EURUSD"
	Side   string `json:"side,omitempty"`   // "long", "short"

	Quantity string `json:"quantity,omitempty"`

	EntryPrice string `json:"entryPrice,omitempty"`
	MarkPrice  string `json:"markPrice,omitempty"`

	UnrealizedPnL string `json:"unrealizedPnl,omitempty"`
	// UnsettledPnL is exchange-specific PnL that includes not-yet-settled components.
	// Some venues (e.g., Orderly) expose "unsettled" PnL that can include realized-but-not-settled PnL.
	UnsettledPnL     string `json:"unsettledPnl,omitempty"`
	LiquidationPrice string `json:"liquidationPrice,omitempty"`

	Leverage int `json:"leverage,omitempty"`

	// IsIsolated indicates isolated vs cross margin when applicable.
	IsIsolated bool `json:"isIsolated,omitempty"`

	InitialMargin     string `json:"initialMargin,omitempty"`
	MaintenanceMargin string `json:"maintenanceMargin,omitempty"`

	// Components provides exchange-specific breakdown values.
	// Common keys: "position_value", "notional", "margin_mode", "funding_rate", "financing_cost", etc.
	Components map[string]string `json:"components,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}
