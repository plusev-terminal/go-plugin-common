package exchange

import (
	"fmt"
	"strings"
	"time"

	acct "github.com/plusev-terminal/go-plugin-common/datasrc/exchange/account"
	tt "github.com/plusev-terminal/go-plugin-common/trading"
	"github.com/plusev-terminal/go-plugin-common/utils"
)

// GetAccountParams contains optional parameters for the getAccount command.
//
// These params are intentionally permissive; when omitted, plugins should return
// their default set of account+scopes.
type GetAccountParams struct {
	Scopes []acct.ScopeType `json:"scopes,omitempty" mapstructure:"scopes"`
}

func (p GetAccountParams) Validate() error { return nil }

// GetBalancesParams contains optional parameters for the getBalances command.
type GetBalancesParams struct {
	Scopes []acct.ScopeType `json:"scopes,omitempty" mapstructure:"scopes"`
}

func (p GetBalancesParams) Validate() error { return nil }

// GetPositionsParams contains optional parameters for the getPositions command.
type GetPositionsParams struct {
	Scopes []acct.ScopeType `json:"scopes,omitempty" mapstructure:"scopes"`
}

func (p GetPositionsParams) Validate() error { return nil }

// GetOrdersParams contains optional parameters for the getOrders command.
//
// When OpenOnly is true, plugins should return only open orders
// (unfilled or partially filled, depending on exchange semantics).
type GetOrdersParams struct {
	OpenOnly bool `json:"openOnly" mapstructure:"openOnly"`
}

func (p GetOrdersParams) Validate() error { return nil }

// CreateOrdersParams contains parameters for the createOrders command.
// Orders supports batch creation; plugins may combine into native bracket/oco orders.
type CreateOrdersParams struct {
	Orders               []CreateOrderParams `json:"orders" mapstructure:"orders" validate:"required,min=1"`
	CancelGroupOnFailure bool                `json:"cancelGroupOnFailure,omitempty" mapstructure:"cancelGroupOnFailure"`
}

// CancelOrdersParams contains parameters for the cancelOrders command.
// If IDs is empty, plugins should cancel all pending orders (and algo orders where supported).
type CancelOrdersParams struct {
	Market tt.Market `json:"market,omitempty" mapstructure:"market"`
	IDs    []string  `json:"ids,omitempty" mapstructure:"ids"`
}

func (p CancelOrdersParams) Validate() error {
	if len(p.IDs) == 0 {
		return nil
	}
	if strings.TrimSpace(p.Market.Symbol) == "" {
		return fmt.Errorf("market.symbol is required")
	}
	return nil
}

func (p CreateOrdersParams) Validate() error {
	if len(p.Orders) == 0 {
		return fmt.Errorf("orders is required")
	}
	for i, order := range p.Orders {
		if order.Market.Symbol == "" {
			return fmt.Errorf("orders[%d].market.symbol is required", i)
		}
		if order.Side == "" || order.Side == acct.OrderSideUnknown {
			return fmt.Errorf("orders[%d].side is required", i)
		}
		if order.Type == "" || order.Type == acct.OrderTypeUnknown {
			return fmt.Errorf("orders[%d].type is required", i)
		}
		if !order.ClosePosition && order.Quantity == "" && order.QuoteQuantity == "" {
			return fmt.Errorf("orders[%d].quantity or orders[%d].quoteQuantity is required", i, i)
		}

		requiresPrice := false
		requiresStop := false
		requiresTakeProfit := false
		requiresTrailing := false

		switch order.Type {
		case acct.OrderTypeLimit:
			requiresPrice = true
		case acct.OrderTypeStopLimit:
			requiresPrice = true
			requiresStop = true
		case acct.OrderTypeTakeProfitLimit:
			requiresPrice = true
			requiresTakeProfit = true
		case acct.OrderTypeStop, acct.OrderTypeStopMarket:
			requiresStop = true
		case acct.OrderTypeTakeProfit, acct.OrderTypeTakeProfitMarket:
			requiresTakeProfit = true
		case acct.OrderTypeTrailingStop:
			requiresTrailing = true
		}

		if requiresPrice && order.Price == "" {
			return fmt.Errorf("orders[%d].price is required for type %s", i, order.Type)
		}
		if requiresStop && order.TriggerPrice == "" && order.StopPrice == "" {
			return fmt.Errorf("orders[%d].triggerPrice or orders[%d].stopPrice is required for type %s", i, i, order.Type)
		}
		if requiresTakeProfit && order.TriggerPrice == "" && order.TakeProfitPrice == "" {
			return fmt.Errorf("orders[%d].triggerPrice or orders[%d].takeProfitPrice is required for type %s", i, i, order.Type)
		}
		if requiresTrailing && order.TrailingDelta == "" {
			return fmt.Errorf("orders[%d].trailingDelta is required for type %s", i, order.Type)
		}
	}
	return nil
}

// CreateOrderParams represents one normalized order request.
// All numeric values are strings to preserve precision.
type CreateOrderParams struct {
	Market tt.Market      `json:"market,omitempty" mapstructure:"market"`
	Scope  acct.ScopeType `json:"scope,omitempty" mapstructure:"scope"`

	Side acct.OrderSide `json:"side" mapstructure:"side" validate:"required"`
	Type acct.OrderType `json:"type" mapstructure:"type" validate:"required"`

	Quantity      string `json:"quantity,omitempty" mapstructure:"quantity"`
	QuoteQuantity string `json:"quoteQuantity,omitempty" mapstructure:"quoteQuantity"`
	Price         string `json:"price,omitempty" mapstructure:"price"`

	TriggerPrice    string `json:"triggerPrice,omitempty" mapstructure:"triggerPrice"`
	StopPrice       string `json:"stopPrice,omitempty" mapstructure:"stopPrice"`
	TakeProfitPrice string `json:"takeProfitPrice,omitempty" mapstructure:"takeProfitPrice"`
	TrailingDelta   string `json:"trailingDelta,omitempty" mapstructure:"trailingDelta"`
	ActivationPrice string `json:"activationPrice,omitempty" mapstructure:"activationPrice"`

	TimeInForce   acct.TimeInForce `json:"timeInForce,omitempty" mapstructure:"timeInForce"`
	ReduceOnly    bool             `json:"reduceOnly,omitempty" mapstructure:"reduceOnly"`
	PostOnly      bool             `json:"postOnly,omitempty" mapstructure:"postOnly"`
	ClosePosition bool             `json:"closePosition,omitempty" mapstructure:"closePosition"`

	ClientOrderID string `json:"clientOrderId,omitempty" mapstructure:"clientOrderId"`
	GroupID       string `json:"groupId,omitempty" mapstructure:"groupId"`
	PositionSide  string `json:"positionSide,omitempty" mapstructure:"positionSide"`

	Extra map[string]any `json:"extra,omitempty" mapstructure:"extra"`
}

func normalizeScopesAny(v any) []acct.ScopeType {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}

	res := make([]acct.ScopeType, 0, len(arr))
	for _, it := range arr {
		s, ok := it.(string)
		if !ok {
			continue
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		res = append(res, acct.ScopeType(s))
	}

	return res
}

// GetAccountParamsFromMap extracts GetAccountParams from a validated map.
func GetAccountParamsFromMap(data map[string]any) GetAccountParams {
	return GetAccountParams{Scopes: normalizeScopesAny(data["scopes"])}
}

// GetBalancesParamsFromMap extracts GetBalancesParams from a validated map.
func GetBalancesParamsFromMap(data map[string]any) GetBalancesParams {
	return GetBalancesParams{Scopes: normalizeScopesAny(data["scopes"])}
}

// GetPositionsParamsFromMap extracts GetPositionsParams from a validated map.
func GetPositionsParamsFromMap(data map[string]any) GetPositionsParams {
	return GetPositionsParams{Scopes: normalizeScopesAny(data["scopes"])}
}

// GetOrdersParamsFromMap extracts GetOrdersParams from a validated map.
func GetOrdersParamsFromMap(data map[string]any) GetOrdersParams {
	openOnly := false
	if v, ok := data["openOnly"].(bool); ok {
		openOnly = v
	}
	return GetOrdersParams{OpenOnly: openOnly}
}

// CreateOrdersParamsFromMap extracts CreateOrdersParams from a validated map.
func CreateOrdersParamsFromMap(data map[string]any) CreateOrdersParams {
	params := CreateOrdersParams{}

	if v, ok := data["cancelGroupOnFailure"].(bool); ok {
		params.CancelGroupOnFailure = v
	}

	if rawOrders, ok := data["orders"].([]any); ok {
		orders := make([]CreateOrderParams, 0, len(rawOrders))

		for _, raw := range rawOrders {
			orderMap, ok := raw.(map[string]any)
			if !ok {
				continue
			}

			var order CreateOrderParams
			if err := utils.MapToStruct(orderMap, &order); err != nil {
				continue
			}

			orders = append(orders, order)
		}

		params.Orders = orders
	}

	return params
}

// CancelOrdersParamsFromMap extracts CancelOrdersParams from a validated map.
func CancelOrdersParamsFromMap(data map[string]any) CancelOrdersParams {
	params := CancelOrdersParams{}
	params.IDs = utils.NormalizeStringsAny(data["ids"])
	if v, ok := data["market"].(map[string]any); ok {
		_ = utils.MapToStruct(v, &params.Market)
	}
	return params
}

// GetMarketParams contains parameters for the getMarket command.
//
// CEX sources typically identify a market by Symbol (+ optional AssetType).
// DEX sources use Address (e.g. pool/contract address) since symbols are ambiguous across pools.
// At least one of Symbol or Address must be provided.
type GetMarketParams struct {
	Symbol    string `json:"symbol,omitempty" mapstructure:"symbol"`
	AssetType string `json:"assetType,omitempty" mapstructure:"assetType"`
	Address   string `json:"address,omitempty" mapstructure:"address"`
}

func (p GetMarketParams) Validate() error {
	if p.Symbol == "" && p.Address == "" {
		return fmt.Errorf("symbol or address is required")
	}
	return nil
}

// GetMarketParamsFromMap extracts GetMarketParams from a validated map.
func GetMarketParamsFromMap(data map[string]any) GetMarketParams {
	return GetMarketParams{
		Symbol:    utils.GetValue[string]("symbol", data),
		AssetType: utils.GetValue[string]("assetType", data),
		Address:   utils.GetValue[string]("address", data),
	}
}

// OHLCVStreamParams contains parameters for the ohlcvStream command
type OHLCVStreamParams struct {
	// Market is required. It provides full context (assetType, base/quote, etc).
	Market    tt.Market `json:"market" mapstructure:"market" validate:"required"`
	Timeframe string    `json:"timeframe" mapstructure:"timeframe" validate:"required"`
}

func (p OHLCVStreamParams) Validate() error {
	if p.Timeframe == "" {
		return fmt.Errorf("timeframe is required")
	}
	if p.Market.Symbol == "" {
		return fmt.Errorf("market.symbol is required")
	}
	return nil
}

// GetOHLCVParams contains parameters for the getOHLCV (historical data) command
type GetOHLCVParams struct {
	Market          tt.Market  `json:"market" mapstructure:"market" validate:"required"`
	Timeframe       string     `json:"timeframe" mapstructure:"timeframe" validate:"required"`
	StartTime       *time.Time `json:"startTime,omitempty" mapstructure:"startTime"`
	EndTime         *time.Time `json:"endTime,omitempty" mapstructure:"endTime"`
	Limit           int        `json:"limit,omitempty" mapstructure:"limit"`
	Sort            string     `json:"sort,omitempty" mapstructure:"sort"`
	CacheForSeconds int        `json:"cacheFor,omitempty" mapstructure:"cacheFor"` // in seconds
}

func (p GetOHLCVParams) Validate() error {
	if p.Timeframe == "" {
		return fmt.Errorf("timeframe is required")
	}

	if p.Market.Symbol == "" {
		return fmt.Errorf("market.symbol is required")
	}
	return nil
}

// OHLCVStreamParamsFromMap extracts OHLCVStreamParams from validated map
func OHLCVStreamParamsFromMap(data map[string]any) OHLCVStreamParams {
	params := OHLCVStreamParams{Timeframe: utils.GetValue[string]("timeframe", data)}
	if v, ok := data["market"].(map[string]any); ok {
		_ = utils.MapToStruct(v, &params.Market)
	}
	return params
}

// GetOHLCVParamsFromMap extracts GetOHLCVParams from validated map
func GetOHLCVParamsFromMap(data map[string]any) GetOHLCVParams {
	params := GetOHLCVParams{
		Timeframe:       utils.GetValue[string]("timeframe", data),
		StartTime:       utils.ExtractTime("startTime", data),
		EndTime:         utils.ExtractTime("endTime", data),
		Limit:           utils.ExtractInt("limit", data),
		Sort:            utils.GetValue[string]("sort", data),
		CacheForSeconds: utils.ExtractInt("cacheFor", data),
	}
	if v, ok := data["market"].(map[string]any); ok {
		_ = utils.MapToStruct(v, &params.Market)
	}
	return params
}

// ---------------------------------------------------------------------------
// Tax / History params
// ---------------------------------------------------------------------------

// GetTradeHistoryParams contains parameters for the getTradeHistory command.
// Plugins paginate through exchange trade history and return []trading.TradeRecord.
type GetTradeHistoryParams struct {
	Market    tt.Market  `json:"market,omitempty" mapstructure:"market"` // Optional symbol filter
	StartTime *time.Time `json:"startTime,omitempty" mapstructure:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty" mapstructure:"endTime"`
	Limit     int        `json:"limit,omitempty" mapstructure:"limit"`   // 0 = exchange default
	Cursor    string     `json:"cursor,omitempty" mapstructure:"cursor"` // Opaque pagination cursor
	Sort      string     `json:"sort,omitempty" mapstructure:"sort"`     // "asc" or "desc"
}

func (p GetTradeHistoryParams) Validate() error { return nil }

// GetTradeHistoryParamsFromMap extracts GetTradeHistoryParams from a validated map.
func GetTradeHistoryParamsFromMap(data map[string]any) GetTradeHistoryParams {
	params := GetTradeHistoryParams{
		StartTime: utils.ExtractTime("startTime", data),
		EndTime:   utils.ExtractTime("endTime", data),
		Limit:     utils.ExtractInt("limit", data),
		Cursor:    utils.GetValue[string]("cursor", data),
		Sort:      utils.GetValue[string]("sort", data),
	}
	if v, ok := data["market"].(map[string]any); ok {
		_ = utils.MapToStruct(v, &params.Market)
	}
	return params
}

// GetOrderHistoryParams contains parameters for the getOrderHistory command.
// Used when an exchange doesn't expose a dedicated trade/fill endpoint — the plugin
// returns order records that the host normalises into TradeRecords.
type GetOrderHistoryParams struct {
	Market    tt.Market  `json:"market,omitempty" mapstructure:"market"`
	StartTime *time.Time `json:"startTime,omitempty" mapstructure:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty" mapstructure:"endTime"`
	Limit     int        `json:"limit,omitempty" mapstructure:"limit"`
	Cursor    string     `json:"cursor,omitempty" mapstructure:"cursor"`
	Sort      string     `json:"sort,omitempty" mapstructure:"sort"`
}

func (p GetOrderHistoryParams) Validate() error { return nil }

// GetOrderHistoryParamsFromMap extracts GetOrderHistoryParams from a validated map.
func GetOrderHistoryParamsFromMap(data map[string]any) GetOrderHistoryParams {
	params := GetOrderHistoryParams{
		StartTime: utils.ExtractTime("startTime", data),
		EndTime:   utils.ExtractTime("endTime", data),
		Limit:     utils.ExtractInt("limit", data),
		Cursor:    utils.GetValue[string]("cursor", data),
		Sort:      utils.GetValue[string]("sort", data),
	}
	if v, ok := data["market"].(map[string]any); ok {
		_ = utils.MapToStruct(v, &params.Market)
	}
	return params
}

// GetTransferHistoryParams contains parameters for the getTransferHistory command.
// Plugins paginate through exchange deposit/withdrawal history and return []trading.TransferRecord.
type GetTransferHistoryParams struct {
	Asset     string     `json:"asset,omitempty" mapstructure:"asset"`   // Optional asset filter
	Action    string     `json:"action,omitempty" mapstructure:"action"` // "deposit", "withdrawal", or "" for both
	StartTime *time.Time `json:"startTime,omitempty" mapstructure:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty" mapstructure:"endTime"`
	Limit     int        `json:"limit,omitempty" mapstructure:"limit"`
	Cursor    string     `json:"cursor,omitempty" mapstructure:"cursor"`
	Sort      string     `json:"sort,omitempty" mapstructure:"sort"`
}

func (p GetTransferHistoryParams) Validate() error { return nil }

// GetTransferHistoryParamsFromMap extracts GetTransferHistoryParams from a validated map.
func GetTransferHistoryParamsFromMap(data map[string]any) GetTransferHistoryParams {
	params := GetTransferHistoryParams{
		Asset:     utils.GetValue[string]("asset", data),
		Action:    utils.GetValue[string]("action", data),
		StartTime: utils.ExtractTime("startTime", data),
		EndTime:   utils.ExtractTime("endTime", data),
		Limit:     utils.ExtractInt("limit", data),
		Cursor:    utils.GetValue[string]("cursor", data),
		Sort:      utils.GetValue[string]("sort", data),
	}
	return params
}
