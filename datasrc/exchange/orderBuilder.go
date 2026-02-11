package exchange

import (
	"fmt"
	"strings"

	acct "github.com/plusev-terminal/go-plugin-common/datasrc/exchange/account"
	tt "github.com/plusev-terminal/go-plugin-common/trading"
)

// ── Standalone constructors ─────────────────────────────────────────────────
// Each constructor only accepts the fields valid for its order type,
// making it structurally impossible to set conflicting fields.

// NewMarketOrder creates a market entry order.
func NewMarketOrder(market tt.Market, side acct.OrderSide, qty string) CreateOrderParams {
	return CreateOrderParams{
		Market:   market,
		Side:     side,
		Type:     acct.OrderTypeMarket,
		Quantity: qty,
	}
}

// NewLimitOrder creates a limit entry order.
func NewLimitOrder(market tt.Market, side acct.OrderSide, qty, price string) CreateOrderParams {
	return CreateOrderParams{
		Market:   market,
		Side:     side,
		Type:     acct.OrderTypeLimit,
		Quantity: qty,
		Price:    price,
	}
}

// WithReduceOnly returns a copy of the order params with ReduceOnly set.
func WithReduceOnly(p CreateOrderParams) CreateOrderParams {
	p.ReduceOnly = true
	return p
}

// NewStopLoss creates a stop-market order intended as a stop loss.
// The order is marked reduce-only.
func NewStopLoss(market tt.Market, side acct.OrderSide, qty, stopPrice string) CreateOrderParams {
	return CreateOrderParams{
		Market:     market,
		Side:       side,
		Type:       acct.OrderTypeStopMarket,
		Quantity:   qty,
		StopPrice:  stopPrice,
		ReduceOnly: true,
	}
}

// NewTakeProfit creates a take-profit-market order.
// The order is marked reduce-only.
func NewTakeProfit(market tt.Market, side acct.OrderSide, qty, tpPrice string) CreateOrderParams {
	return CreateOrderParams{
		Market:          market,
		Side:            side,
		Type:            acct.OrderTypeTakeProfitMarket,
		Quantity:        qty,
		TakeProfitPrice: tpPrice,
		ReduceOnly:      true,
	}
}

// NewTrailingStop creates a trailing-stop order.
// The order is marked reduce-only.
func NewTrailingStop(market tt.Market, side acct.OrderSide, qty, trailingDelta string) CreateOrderParams {
	return CreateOrderParams{
		Market:        market,
		Side:          side,
		Type:          acct.OrderTypeTrailingStop,
		Quantity:      qty,
		TrailingDelta: trailingDelta,
		ReduceOnly:    true,
	}
}

// ── Standalone order builder ────────────────────────────────────────────────
// Fluent builder for a single CreateOrderParams with chainable setters.
// Also supports bracket orders (entry + SL/TP) via BuildGroup().
//
// Usage (single order):
//
//	order := exchange.NewOrder(market, acct.OrderSideBuy).
//	    Market("1.5").
//	    SetReduceOnly().
//	    ClientOrderID("my-id").
//	    Build()
//
// Usage (bracket group):
//
//	orders, err := exchange.NewOrder(market, acct.OrderSideBuy).
//	    Market("1.5").
//	    WithStopLoss("1800.00").
//	    WithTakeProfit("2200.00").
//	    GroupID("glf_123_456").
//	    BuildGroup()

// OrderBuilder constructs a single CreateOrderParams via method chaining.
// For bracket orders, attach protection orders with WithStopLoss / WithTakeProfit
// and call BuildGroup instead of Build.
type OrderBuilder struct {
	p  CreateOrderParams
	sl *CreateOrderParams
	tp *CreateOrderParams
}

// NewOrder starts building a standalone order for the given market and side.
func NewOrder(market tt.Market, side acct.OrderSide) *OrderBuilder {
	return &OrderBuilder{p: CreateOrderParams{Market: market, Side: side}}
}

// Market sets the order type to market with the given base quantity.
func (b *OrderBuilder) Market(qty string) *OrderBuilder {
	b.p.Type = acct.OrderTypeMarket
	b.p.Quantity = qty
	return b
}

// Limit sets the order type to limit with the given quantity and price.
func (b *OrderBuilder) Limit(qty, price string) *OrderBuilder {
	b.p.Type = acct.OrderTypeLimit
	b.p.Quantity = qty
	b.p.Price = price
	return b
}

// StopMarket sets the order type to stop-market with the given quantity and stop price.
func (b *OrderBuilder) StopMarket(qty, stopPrice string) *OrderBuilder {
	b.p.Type = acct.OrderTypeStopMarket
	b.p.Quantity = qty
	b.p.StopPrice = stopPrice
	return b
}

// StopLimit sets the order type to stop-limit with the given quantity, limit price, and stop price.
func (b *OrderBuilder) StopLimit(qty, price, stopPrice string) *OrderBuilder {
	b.p.Type = acct.OrderTypeStopLimit
	b.p.Quantity = qty
	b.p.Price = price
	b.p.StopPrice = stopPrice
	return b
}

// TakeProfitMarket sets the order type to take-profit-market with the given quantity and TP price.
func (b *OrderBuilder) TakeProfitMarket(qty, tpPrice string) *OrderBuilder {
	b.p.Type = acct.OrderTypeTakeProfitMarket
	b.p.Quantity = qty
	b.p.TakeProfitPrice = tpPrice
	return b
}

// TrailingStop sets the order type to trailing-stop with the given quantity and trailing delta.
func (b *OrderBuilder) TrailingStop(qty, trailingDelta string) *OrderBuilder {
	b.p.Type = acct.OrderTypeTrailingStop
	b.p.Quantity = qty
	b.p.TrailingDelta = trailingDelta
	return b
}

// SetReduceOnly marks the order as reduce-only.
func (b *OrderBuilder) SetReduceOnly() *OrderBuilder {
	b.p.ReduceOnly = true
	return b
}

// PostOnly marks the order as post-only.
func (b *OrderBuilder) PostOnly() *OrderBuilder {
	b.p.PostOnly = true
	return b
}

// ClosePosition marks the order as close-position.
func (b *OrderBuilder) ClosePosition() *OrderBuilder {
	b.p.ClosePosition = true
	return b
}

// TimeInForce sets the time-in-force policy.
func (b *OrderBuilder) TimeInForce(tif acct.TimeInForce) *OrderBuilder {
	b.p.TimeInForce = tif
	return b
}

// ClientOrderID sets a client-assigned order ID.
func (b *OrderBuilder) ClientOrderID(id string) *OrderBuilder {
	b.p.ClientOrderID = id
	return b
}

// GroupID assigns the order to a group.
func (b *OrderBuilder) GroupID(id string) *OrderBuilder {
	b.p.GroupID = id
	return b
}

// PositionSide sets the position side (for hedge-mode exchanges).
func (b *OrderBuilder) PositionSide(side string) *OrderBuilder {
	b.p.PositionSide = side
	return b
}

// QuoteQuantity sets the quote-currency quantity (e.g. spend $100 worth).
func (b *OrderBuilder) QuoteQuantity(qty string) *OrderBuilder {
	b.p.QuoteQuantity = qty
	return b
}

// ActivationPrice sets the activation price (used by some trailing stop implementations).
func (b *OrderBuilder) ActivationPrice(price string) *OrderBuilder {
	b.p.ActivationPrice = price
	return b
}

// Extra merges additional key-value pairs into the order's Extra map.
func (b *OrderBuilder) Extra(kv map[string]any) *OrderBuilder {
	if b.p.Extra == nil {
		b.p.Extra = make(map[string]any, len(kv))
	}
	for k, v := range kv {
		b.p.Extra[k] = v
	}
	return b
}

// WithStopLoss attaches a stop-loss protection order for use with BuildGroup.
// Side and quantity are inferred from the entry order.
func (b *OrderBuilder) WithStopLoss(stopPrice string) *OrderBuilder {
	if isEmptyPrice(stopPrice) {
		return b
	}
	b.sl = &CreateOrderParams{
		Type:      acct.OrderTypeStopMarket,
		StopPrice: stopPrice,
	}
	return b
}

// WithTakeProfit attaches a take-profit protection order for use with BuildGroup.
// Side and quantity are inferred from the entry order.
func (b *OrderBuilder) WithTakeProfit(tpPrice string) *OrderBuilder {
	if isEmptyPrice(tpPrice) {
		return b
	}
	b.tp = &CreateOrderParams{
		Type:            acct.OrderTypeTakeProfitMarket,
		TakeProfitPrice: tpPrice,
	}
	return b
}

// Build returns the final CreateOrderParams for a single order.
func (b *OrderBuilder) Build() CreateOrderParams {
	return b.p
}

// BuildGroup validates and returns a bracket order group (entry + SL and/or TP).
// The entry's side and quantity are propagated to protection orders (with the
// opposite side), and all orders share the configured GroupID.
func (b *OrderBuilder) BuildGroup() ([]CreateOrderParams, error) {
	if b.p.Type == "" {
		return nil, fmt.Errorf("order group requires an entry order type (call Market, Limit, etc.)")
	}

	if b.sl == nil && b.tp == nil {
		return nil, fmt.Errorf("order group requires at least one protection order (call WithStopLoss or WithTakeProfit)")
	}

	closeSide := oppositeSide(b.p.Side)
	orders := []CreateOrderParams{b.p}

	if b.sl != nil {
		b.sl.Market = b.p.Market
		b.sl.Side = closeSide
		b.sl.Quantity = b.p.Quantity
		b.sl.ReduceOnly = true
		b.sl.GroupID = b.p.GroupID
		orders = append(orders, *b.sl)
	}

	if b.tp != nil {
		b.tp.Market = b.p.Market
		b.tp.Side = closeSide
		b.tp.Quantity = b.p.Quantity
		b.tp.ReduceOnly = true
		b.tp.GroupID = b.p.GroupID
		orders = append(orders, *b.tp)
	}

	return orders, nil
}

// ── Deprecated ──────────────────────────────────────────────────────────────
// NewOrderGroup is deprecated. Use NewOrder(...).BuildGroup() instead.
func NewOrderGroup(market tt.Market, groupID string) *OrderBuilder {
	return &OrderBuilder{p: CreateOrderParams{Market: market, GroupID: groupID}}
}

// MarketEntry is a convenience alias used by the deprecated NewOrderGroup flow.
// Prefer NewOrder(market, side).Market(qty).GroupID(id) instead.
func (b *OrderBuilder) MarketEntry(side acct.OrderSide, qty string) *OrderBuilder {
	b.p.Side = side
	b.p.Type = acct.OrderTypeMarket
	b.p.Quantity = qty
	return b
}

// isEmptyPrice returns true for blank strings and "NaN" values.
func isEmptyPrice(v string) bool {
	v = strings.TrimSpace(v)
	return v == "" || strings.EqualFold(v, "nan")
}

func oppositeSide(side acct.OrderSide) acct.OrderSide {
	switch side {
	case acct.OrderSideBuy:
		return acct.OrderSideSell
	case acct.OrderSideSell:
		return acct.OrderSideBuy
	default:
		return acct.OrderSideUnknown
	}
}
