package account

import (
	"encoding/json"
	"strings"
	"time"
)

// GetOrdersRequest are the optional params for the getOrders command.
// Implementations may use Scopes to avoid unnecessary API calls.
type GetOrdersRequest struct {
	ScopeFilter

	// OpenOnly requests only open orders (unfilled or partially filled).
	OpenOnly bool `json:"open_only,omitempty"`
}

// OrdersResponse is the response type for the getOrders command.
// Orders are grouped by scope so consumers can display/aggregate without having to
// infer scope from symbol naming.
type OrdersResponse struct {
	FetchedAt time.Time `json:"fetched_at,omitempty"`

	Scopes map[ScopeType]OrderScope `json:"scopes,omitempty"`

	// Raw stores the unmodified exchange payload for debugging/support.
	Raw json.RawMessage `json:"raw,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// OrderScope represents orders within a specific scope.
type OrderScope struct {
	ScopeID string `json:"scope_id,omitempty"`

	// Network identifies the chain/network when ScopeType is "wallet" (self-custody).
	// Examples: "ethereum", "arbitrum", "solana", "bitcoin".
	Network Network `json:"network,omitempty"`

	// ChainID is the numeric EVM chain id when applicable (e.g., 1 for Ethereum mainnet).
	// Leave as 0 when not applicable or unknown.
	ChainID int `json:"chain_id,omitempty"`

	Orders []Order `json:"orders,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// OrderSide is a normalized order side.
type OrderSide string

func (o OrderSide) Invert() OrderSide {
	switch o {
	case OrderSideBuy:
		return OrderSideSell
	case OrderSideSell:
		return OrderSideBuy
	default:
		return OrderSideUnknown
	}
}

const (
	OrderSideUnknown OrderSide = "unknown"
	OrderSideBuy     OrderSide = "buy"
	OrderSideSell    OrderSide = "sell"
)

// OrderType is a normalized order type.
type OrderType string

const (
	OrderTypeUnknown          OrderType = "unknown"
	OrderTypeLimit            OrderType = "limit"
	OrderTypeMarket           OrderType = "market"
	OrderTypeStop             OrderType = "stop"
	OrderTypeStopLimit        OrderType = "stop_limit"
	OrderTypeStopMarket       OrderType = "stop_market"
	OrderTypeTakeProfit       OrderType = "take_profit"
	OrderTypeTakeProfitLimit  OrderType = "take_profit_limit"
	OrderTypeTakeProfitMarket OrderType = "take_profit_market"
	OrderTypeTrailingStop     OrderType = "trailing_stop"
)

// OrderStatus is a normalized order status.
type OrderStatus string

const (
	OrderStatusUnknown         OrderStatus = "unknown"
	OrderStatusNew             OrderStatus = "new"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
	OrderStatusFilled          OrderStatus = "filled"
	OrderStatusCanceled        OrderStatus = "canceled"
	OrderStatusRejected        OrderStatus = "rejected"
	OrderStatusExpired         OrderStatus = "expired"
	OrderStatusIncomplete      OrderStatus = "incomplete"
	OrderStatusCompleted       OrderStatus = "completed"
)

// TimeInForce is a normalized time-in-force policy.
type TimeInForce string

const (
	TimeInForceUnknown  TimeInForce = "unknown"
	TimeInForceGTC      TimeInForce = "gtc"
	TimeInForceIOC      TimeInForce = "ioc"
	TimeInForceFOK      TimeInForce = "fok"
	TimeInForcePostOnly TimeInForce = "post_only"
	TimeInForceDay      TimeInForce = "day"
	TimeInForceGTX      TimeInForce = "gtx"
	TimeInForceGTD      TimeInForce = "gtd"
)

// TriggerSource selects the price a conditional order triggers against.
// Binance calls it workingType, Bybit triggerBy, Orderly trigger_price_type —
// one concept.
//
// The empty TriggerSource means "the venue's default", which differs per
// venue (Binance and Bybit trigger on the last price, Orderly on the mark
// price); venues publish theirs as OrderCapabilities.DefaultTriggerSource.
// Consumers that care which price trips a stop must set it explicitly.
type TriggerSource string

const (
	// TriggerSourceLast is the last traded price.
	TriggerSourceLast TriggerSource = "last"
	// TriggerSourceMark is the mark price.
	TriggerSourceMark TriggerSource = "mark"
	// TriggerSourceIndex is the index price.
	TriggerSourceIndex TriggerSource = "index"
)

// Valid reports whether s is one of the defined trigger sources. The empty
// TriggerSource (venue default) is not Valid — check it separately.
func (s TriggerSource) Valid() bool {
	switch s {
	case TriggerSourceLast, TriggerSourceMark, TriggerSourceIndex:
		return true
	}

	return false
}

// OptionalPrice normalizes a venue price field for the normalized Order and
// Position shapes: an absent price is the empty string, never a "0"
// placeholder. Venues fill unused price fields with zeros (Binance sends
// price "0" on stop-market orders and stopPrice "0" on plain limit orders),
// and a consumer must not read those as real prices.
func OptionalPrice(s string) string {
	s = strings.TrimSpace(s)

	if s == "" || strings.Trim(s, "0.+-") == "" {
		return ""
	}

	return s
}

// Order represents a normalized order across exchanges.
// All numeric values are strings to preserve precision.
type Order struct {
	OrderID       string `json:"order_id,omitempty"`
	ClientOrderID string `json:"client_order_id,omitempty"`

	Symbol string      `json:"symbol,omitempty"`
	Side   OrderSide   `json:"side,omitempty"`
	Type   OrderType   `json:"type,omitempty"`
	Status OrderStatus `json:"status,omitempty"`

	// Price is the limit price; empty for market-type orders. Absent prices
	// (here and in TriggerPrice/AveragePrice) are "", never "0" — plugins
	// normalize venue placeholders with OptionalPrice.
	Price             string `json:"price,omitempty"`
	Quantity          string `json:"quantity,omitempty"`
	ExecutedQuantity  string `json:"executed_quantity,omitempty"`
	RemainingQuantity string `json:"remaining_quantity,omitempty"`
	AveragePrice      string `json:"average_price,omitempty"`

	// TriggerPrice is the activation price of a conditional order (stop,
	// take-profit, trailing activation); empty for non-conditional orders.
	// Plugins populate it from whichever native field the order type
	// consumes; it is THE field consumers read — do not dig through
	// Components for trigger prices.
	TriggerPrice string `json:"trigger_price,omitempty"`
	// TriggerSource is the price the conditional order triggers against
	// (last/mark/index); empty means the venue default.
	TriggerSource TriggerSource `json:"trigger_source,omitempty"`

	TimeInForce TimeInForce `json:"time_in_force,omitempty"`
	ReduceOnly  bool        `json:"reduce_only,omitempty"`
	PostOnly    bool        `json:"post_only,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	// Components provides exchange-specific breakdown values beyond the
	// first-class fields (e.g. "fee", "fee_asset", venue-native order ids).
	// Trigger prices live in TriggerPrice — not here.
	Components map[string]string `json:"components,omitempty"`
	Extra      map[string]any    `json:"extra,omitempty"`
}

// CreateOrdersResponse is the response type for the createOrders command.
// Results are returned per order to support batch requests.
type CreateOrdersResponse struct {
	SubmittedAt time.Time `json:"submitted_at,omitempty"`

	Results []CreateOrderResult `json:"results,omitempty"`

	// Raw stores the unmodified exchange payload for debugging/support.
	Raw json.RawMessage `json:"raw,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// CreateOrderResult contains the per-order result for createOrders.
type CreateOrderResult struct {
	Index int `json:"index,omitempty"`

	ClientOrderID string `json:"client_order_id,omitempty"`
	GroupID       string `json:"group_id,omitempty"`

	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	Order *Order `json:"order,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// CancelOrdersResponse is the response type for the cancelOrders command.
type CancelOrdersResponse struct {
	SubmittedAt time.Time `json:"submitted_at,omitempty"`

	Results []CancelOrderResult `json:"results,omitempty"`

	// Raw stores the unmodified exchange payload for debugging/support.
	Raw json.RawMessage `json:"raw,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// CancelOrderResult contains the per-order result for cancelOrders.
type CancelOrderResult struct {
	Index int `json:"index,omitempty"`

	OrderID       string `json:"order_id,omitempty"`
	ClientOrderID string `json:"client_order_id,omitempty"`

	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// ModifyOrdersResponse is the response type for the modifyOrders command.
// Results are returned per order to support batch requests.
type ModifyOrdersResponse struct {
	SubmittedAt time.Time `json:"submitted_at,omitempty"`

	Results []ModifyOrderResult `json:"results,omitempty"`

	// Raw stores the unmodified exchange payload for debugging/support.
	Raw json.RawMessage `json:"raw,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}

// ModifyOrderResult contains the per-order result for modifyOrders.
//
// The order ID is stable on a native modify; when the exchange (or the
// plugin, via a cancelReplace) had to replace the order, Order carries the
// new order and its OrderID differs from the requested one.
type ModifyOrderResult struct {
	Index int `json:"index,omitempty"`

	OrderID string `json:"order_id,omitempty"`

	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	Order *Order `json:"order,omitempty"`

	Extra map[string]any `json:"extra,omitempty"`
}
