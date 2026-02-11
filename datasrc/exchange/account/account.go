package account

import (
	"encoding/json"
	"time"
)

// AccountType represents the type of exchange account
type AccountType string

const (
	AccountTypeSpot    AccountType = "spot"
	AccountTypeMargin  AccountType = "margin"
	AccountTypeFutures AccountType = "futures"
	AccountTypeFunding AccountType = "funding"
	AccountTypeEarn    AccountType = "earn"
	AccountTypeUnknown AccountType = "unknown"
)

// CustodyModel describes who controls the assets (or whether they are synthetic).
// This is independent from "spot/margin/futures" scopes.
type CustodyModel string

const (
	CustodyModelUnknown         CustodyModel = "unknown"
	CustodyModelExchangeCustody CustodyModel = "exchange_custody" // CEX: venue controls keys
	CustodyModelSelfCustody     CustodyModel = "self_custody"     // DEX / wallet-connected: user controls keys
	CustodyModelSynthetic       CustodyModel = "synthetic"        // CFDs / IOUs: no on-chain/withdrawable asset ownership
)

// Account represents a normalized structure for user account details across different exchanges.
// All monetary values are strings to preserve precision; parse with a decimal library when needed.
type Account struct {
	Exchange   string      `json:"exchange,omitempty"`   // "binance", "orderly", "woox", etc.
	AccountID  string      `json:"account_id,omitempty"` // exchange account id
	Subaccount string      `json:"subaccount,omitempty"` // if applicable
	Type       AccountType `json:"type,omitempty"`
	Status     string      `json:"status,omitempty"` // "active", "locked", "verified", etc.

	// BaseCurrency is the account's base/settlement currency when applicable.
	// Common for legacy brokers (e.g., "USD", "EUR"); may be empty for exchanges.
	BaseCurrency string `json:"base_currency,omitempty"`

	// CustodyModel indicates the custody model for this account.
	// Examples:
	// - Exchange custodial spot (Binance): "exchange_custody"
	// - Self-custodial DEX wallet: "self_custody"
	// - Legacy CFD broker: "synthetic"
	CustodyModel CustodyModel `json:"custody_model,omitempty"`

	// IsCustodial indicates whether the account represents custody of real assets.
	// For CFDs/synthetic brokers this should be false; for spot (custodial or self-custody) it is typically true.
	// Prefer using CustodyModel when you need to distinguish custodial vs self-custody.
	IsCustodial bool `json:"is_custodial,omitempty"`

	// Features is a set of capability tags (best effort), e.g. "swap_free", "commission_based".
	Features []string `json:"features,omitempty"`

	// MaxLeverage is a best-effort account-level maximum leverage (mainly for derivatives).
	// Not all exchanges provide this; when unknown it will be 0.
	MaxLeverage int `json:"max_leverage,omitempty"`

	FetchedAt time.Time `json:"fetched_at,omitempty"`

	// Scopes partitions balances by wallet/category (spot, margin, futures, etc.)
	// This avoids mixing incompatible balance buckets.
	Scopes map[ScopeType]BalanceScope `json:"scopes,omitempty"`

	// Raw stores the unmodified exchange payload for debugging/support.
	Raw json.RawMessage `json:"raw,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}
