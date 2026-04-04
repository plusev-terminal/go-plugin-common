package account

import (
	"encoding/json"
	"time"
)

// BalancesResponse is the response type for the getBalances command.
// It intentionally contains only balance-related data, plus metadata that applies only
// to the balances snapshot (e.g., FetchedAt), to avoid confusion with full account details.
type BalancesResponse struct {
	FetchedAt time.Time `json:"fetched_at,omitempty"`

	// Scopes partitions balances by wallet/category (spot, margin, futures, etc.).
	Scopes map[ScopeType]BalanceScope `json:"scopes,omitempty"`

	// Raw stores the unmodified exchange payload for debugging/support.
	Raw json.RawMessage `json:"raw,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// ScopeFilter controls which scopes the plugin should fetch.
// When empty/nil, the plugin should use its default behaviour (usually fetch what it can cheaply).
// When provided, the plugin should try to avoid calling APIs outside the requested scopes.
type ScopeFilter struct {
	Scopes []ScopeType `json:"scopes,omitempty"`
}

// GetAccountRequest are the optional params for the getAccount command.
// Implementations may use Scopes to avoid unnecessary API calls.
type GetAccountRequest struct {
	ScopeFilter
}

// GetBalancesRequest are the optional params for the getBalances command.
// Implementations may use Scopes to avoid unnecessary API calls.
type GetBalancesRequest struct {
	ScopeFilter
}

// ScopeType represents different wallet/balance categories.
// These are used as keys in Account.Scopes and BalancesResponse.Scopes.
type ScopeType string

const (
	ScopeSpot ScopeType = "spot"
	// ScopeWallet represents a self-custody wallet scope (e.g., EVM address, Solana pubkey).
	// Use this for DEX / wallet-connected integrations.
	ScopeWallet         ScopeType = "wallet"
	ScopeMargin         ScopeType = "margin"
	ScopeMarginIsolated ScopeType = "margin_isolated"
	ScopeFutures        ScopeType = "futures"
	// ScopeCollateral represents assets posted/available as collateral for trading.
	// For venues where users deposit assets and those assets back margin (e.g., perps),
	// map holdings/collateral balances here.
	ScopeCollateral ScopeType = "collateral"
	ScopeEarn       ScopeType = "earn"
	// ScopeCFD represents a unified margin pool common in CFD/legacy brokers.
	ScopeCFD ScopeType = "cfd"
	// ScopeUnified is a generic single-pool scope when the broker/exchange doesn't expose subaccounts.
	ScopeUnified ScopeType = "unified"
)

// Network identifies a chain/network for self-custody (and on-chain protocols).
// It is a string so plugins can still use custom networks.
type Network string

const (
	NetworkUnknown   Network = "unknown"
	NetworkEthereum  Network = "ethereum"
	NetworkArbitrum  Network = "arbitrum"
	NetworkOptimism  Network = "optimism"
	NetworkBase      Network = "base"
	NetworkPolygon   Network = "polygon"
	NetworkAvalanche Network = "avalanche"
	NetworkBSC       Network = "bsc"
	NetworkSolana    Network = "solana"
	NetworkBitcoin   Network = "bitcoin"
)

// SpotBalance is a convenience accessor for spot wallet balances.
func (a *Account) SpotBalance(asset string) (AssetBalance, bool) {
	if a == nil {
		return AssetBalance{}, false
	}
	if scope, ok := a.Scopes[ScopeSpot]; ok {
		bal, ok := scope.Balances[asset]
		return bal, ok
	}
	return AssetBalance{}, false
}

// BalanceScope represents balances within a specific wallet/category.
type BalanceScope struct {
	ScopeID string `json:"scope_id,omitempty"` // walletId, portfolioId if exchange provides it

	// Network identifies the chain/network when ScopeType is "wallet" (self-custody).
	// Examples: "ethereum", "arbitrum", "solana", "bitcoin".
	Network Network `json:"network,omitempty"`

	// ChainID is the numeric EVM chain id when applicable (e.g., 1 for Ethereum mainnet).
	// Leave as 0 when not applicable or unknown.
	ChainID int `json:"chain_id,omitempty"`

	// Features is a set of capability tags specific to this scope (best effort).
	// Examples: "cross", "isolated", "swap_free", "commission_based".
	Features []string `json:"features,omitempty"`

	// Balances maps asset symbol (e.g., "BTC", "USDT") to balance details.
	Balances map[string]AssetBalance `json:"balances,omitempty"`

	// State holds scope-level account/risk state (optional), such as margin availability.
	State *ScopeState `json:"state,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// AssetBalance represents the normalized balance details for a specific asset.
type AssetBalance struct {
	Asset string `json:"asset,omitempty"`

	// Core balance fields (all values as strings for decimal precision)
	Total               string `json:"total,omitempty"`
	AvailableToTrade    string `json:"available_to_trade,omitempty"`
	AvailableToWithdraw string `json:"available_to_withdraw,omitempty"`

	// Components provides a flexible breakdown for exchange-specific buckets.
	// Common keys: "open_orders", "staking_reward", "collateral", "pending_withdrawal",
	// "unrealized_pnl", "initial_margin", "earn", "locked", "frozen", etc.
	Components map[string]string `json:"components,omitempty"`

	// Margin-specific fields (when applicable)
	Borrowed string `json:"borrowed,omitempty"`
	Interest string `json:"interest,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// Locked returns a best-effort locked/held amount from components.
// Prefer reading explicit component keys rather than deriving from Total/Available,
// since exchanges don't share consistent semantics.
func (b AssetBalance) Locked() string {
	if v, ok := b.Components["locked"]; ok {
		return v
	}
	if v, ok := b.Components["open_orders"]; ok {
		return v
	}
	return ""
}

// ScopeState holds scope-level aggregates (margin/risk), when the venue provides them.
// All numeric values are strings.
type ScopeState struct {
	Equity            string `json:"equity,omitempty"`
	UsedMargin        string `json:"used_margin,omitempty"`
	AvailableMargin   string `json:"available_margin,omitempty"`
	MaintenanceMargin string `json:"maintenance_margin,omitempty"`
	InitialMargin     string `json:"initial_margin,omitempty"`
	UnrealizedPnL     string `json:"unrealized_pnl,omitempty"`
	MarginRatio       string `json:"margin_ratio,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}
