package exchange

import (
	"fmt"
	"strings"

	acct "github.com/plusev-terminal/go-plugin-common/v2/datasrc/exchange/account"
	tt "github.com/plusev-terminal/go-plugin-common/v2/trading"
)

// OrderCapabilities declares which order features a venue supports.
//
// Every plugin that registers createOrders publishes one (plugin.Plugin
// implementations add GetOrderCapabilities — see plugin.OrderCapabilityProvider).
// The host checks every order against it BEFORE dispatch and rejects any
// feature the plugin did not declare, so an order never reaches a venue that
// would silently ignore part of it. The check runs on the host's (newest)
// contract, which is what makes it forward compatible: a field added to
// CreateOrderParams after a plugin was built is undeclared for that plugin
// and therefore rejected — it can never be dropped unnoticed.
//
// Everything is opt-in: a zero value supports nothing. Empty optional order
// fields always pass (they mean "venue default").
type OrderCapabilities struct {
	// AssetTypes the venue accepts orders for (e.g. spot, perpetual).
	AssetTypes []tt.AssetType `json:"assetTypes"`
	// OrderTypes the venue can place.
	OrderTypes []acct.OrderType `json:"orderTypes"`
	// TimeInForce values the venue honors (an empty TimeInForce always passes).
	TimeInForce []acct.TimeInForce `json:"timeInForce,omitempty"`

	// TriggerSources the venue honors for conditional orders, and the one it
	// uses when an order leaves TriggerSource empty.
	TriggerSources       []acct.TriggerSource `json:"triggerSources,omitempty"`
	DefaultTriggerSource acct.TriggerSource   `json:"defaultTriggerSource,omitempty"`

	// Trailing-stop distance units and activation price.
	TrailingAmount  bool `json:"trailingAmount,omitempty"`
	TrailingPercent bool `json:"trailingPercent,omitempty"`
	ActivationPrice bool `json:"activationPrice,omitempty"`

	ReduceOnly    bool `json:"reduceOnly,omitempty"`
	PostOnly      bool `json:"postOnly,omitempty"`
	ClosePosition bool `json:"closePosition,omitempty"`
	QuoteQuantity bool `json:"quoteQuantity,omitempty"`
	ClientOrderID bool `json:"clientOrderId,omitempty"`
	// PositionSide marks hedge-mode support (long/short legs per symbol).
	PositionSide bool `json:"positionSide,omitempty"`

	// CancelGroupOnFailure marks support for the batch flag of the same name:
	// when one order of a group fails, the plugin cancels the group's orders
	// that were already placed (no entry survives without its protection).
	CancelGroupOnFailure bool `json:"cancelGroupOnFailure,omitempty"`

	// ExtraKeys lists the venue-specific CreateOrderParams.Extra keys the
	// plugin reads; any other Extra key is rejected.
	ExtraKeys []string `json:"extraKeys,omitempty"`

	// legacy marks the implicit baseline for plugins built before v2.1.0
	// (see LegacyOrderCapabilities). Never serialized.
	legacy bool
}

// LegacyOrderCapabilities is the baseline the host assumes for plugins built
// before go-plugin-common v2.1.0, which cannot declare capabilities: the v1
// order surface they were written against stays allowed, while everything
// the contract added later (trailing units, trigger source, GTD expiry,
// group cancellation) is rejected — those binaries do not know the fields
// and would drop them. Asset types and Extra keys are not checked (v1 had no
// way to declare them).
func LegacyOrderCapabilities() *OrderCapabilities {
	return &OrderCapabilities{
		OrderTypes: []acct.OrderType{
			acct.OrderTypeLimit, acct.OrderTypeMarket,
			acct.OrderTypeStop, acct.OrderTypeStopLimit, acct.OrderTypeStopMarket,
			acct.OrderTypeTakeProfit, acct.OrderTypeTakeProfitLimit, acct.OrderTypeTakeProfitMarket,
		},
		TimeInForce: []acct.TimeInForce{
			acct.TimeInForceGTC, acct.TimeInForceIOC, acct.TimeInForceFOK,
			acct.TimeInForcePostOnly, acct.TimeInForceDay, acct.TimeInForceGTX,
		},
		ActivationPrice: true,
		ReduceOnly:      true,
		PostOnly:        true,
		ClosePosition:   true,
		QuoteQuantity:   true,
		ClientOrderID:   true,
		PositionSide:    true,
		legacy:          true,
	}
}

// Legacy reports whether c is the implicit pre-v2.1.0 baseline.
func (c *OrderCapabilities) Legacy() bool {
	return c != nil && c.legacy
}

// SupportsAssetType reports whether orders on markets of type a are accepted.
func (c *OrderCapabilities) SupportsAssetType(a tt.AssetType) bool {
	if c == nil {
		return false
	}

	return c.legacy || contains(c.AssetTypes, a)
}

// SupportsOrderType reports whether the venue can place orders of type t.
func (c *OrderCapabilities) SupportsOrderType(t acct.OrderType) bool {
	return c != nil && contains(c.OrderTypes, t)
}

// Check returns an error naming the first feature o uses that c does not
// declare. It does not re-run CreateOrdersParams.Validate — call both.
func (c *OrderCapabilities) Check(o CreateOrderParams) error {
	if c == nil {
		return fmt.Errorf("the venue declares no order capabilities")
	}

	if !c.legacy && !contains(c.AssetTypes, o.Market.AssetType) {
		return fmt.Errorf("%s markets are not tradable on this venue", describeAssetType(o.Market.AssetType))
	}

	if !contains(c.OrderTypes, o.Type) {
		return fmt.Errorf("order type %s is not supported by this venue", o.Type)
	}

	if o.TimeInForce != "" && !contains(c.TimeInForce, o.TimeInForce) {
		return fmt.Errorf("timeInForce %s is not supported by this venue", o.TimeInForce)
	}

	if o.TriggerSource != "" && o.TriggerSource != c.DefaultTriggerSource && !contains(c.TriggerSources, o.TriggerSource) {
		return fmt.Errorf("triggerSource %s is not supported by this venue", o.TriggerSource)
	}

	flags := []struct {
		used      bool
		supported bool
		name      string
	}{
		{o.TrailingAmount != "", c.TrailingAmount, "trailingAmount"},
		{o.TrailingPercent != "", c.TrailingPercent, "trailingPercent"},
		{o.ActivationPrice != "", c.ActivationPrice, "activationPrice"},
		{o.ReduceOnly, c.ReduceOnly, "reduceOnly"},
		{o.PostOnly, c.PostOnly, "postOnly"},
		{o.ClosePosition, c.ClosePosition, "closePosition"},
		{o.QuoteQuantity != "", c.QuoteQuantity, "quoteQuantity"},
		{o.ClientOrderID != "", c.ClientOrderID, "clientOrderId"},
		{o.PositionSide != "", c.PositionSide, "positionSide"},
	}

	for _, f := range flags {
		if f.used && !f.supported {
			return fmt.Errorf("%s is not supported by this venue", f.name)
		}
	}

	if !c.legacy {
		for key := range o.Extra {
			if !contains(c.ExtraKeys, key) {
				return fmt.Errorf("extra.%s is not supported by this venue", key)
			}
		}
	}

	return nil
}

// CheckAll checks every order of a batch and the batch-level flags.
func (c *OrderCapabilities) CheckAll(p CreateOrdersParams) error {
	if p.CancelGroupOnFailure && (c == nil || !c.CancelGroupOnFailure) {
		return fmt.Errorf("cancelGroupOnFailure is not supported by this venue")
	}

	for i, o := range p.Orders {
		if err := c.Check(o); err != nil {
			return fmt.Errorf("orders[%d]: %w", i, err)
		}
	}

	return nil
}

func describeAssetType(a tt.AssetType) string {
	if a == "" {
		return "unspecified-type"
	}

	return string(a)
}

func contains[T comparable](list []T, v T) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}

	return false
}

// joinStrings is a tiny helper for capability summaries.
func joinStrings[T ~string](list []T) string {
	parts := make([]string, 0, len(list))

	for _, item := range list {
		parts = append(parts, string(item))
	}

	return strings.Join(parts, ", ")
}

// String summarizes the declaration for logs.
func (c *OrderCapabilities) String() string {
	if c == nil {
		return "none"
	}

	if c.legacy {
		return "legacy (pre-v2.1.0 baseline)"
	}

	return fmt.Sprintf("assetTypes=[%s] orderTypes=[%s] tif=[%s] triggerSources=[%s] default=%s",
		joinStrings(c.AssetTypes), joinStrings(c.OrderTypes), joinStrings(c.TimeInForce),
		joinStrings(c.TriggerSources), c.DefaultTriggerSource)
}
