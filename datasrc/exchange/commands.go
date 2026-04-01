package exchange

const (
	CMD_GET_ACCOUNT    = "getAccount"
	CMD_GET_BALANCES   = "getBalances"
	CMD_GET_POSITIONS  = "getPositions"
	CMD_GET_ORDERS     = "getOrders"
	CMD_CREATE_ORDERS  = "createOrders"
	CMD_CANCEL_ORDERS  = "cancelOrders"
	CMD_GET_MARKETS    = "getMarkets"
	CMD_GET_MARKET     = "getMarket"
	CMD_GET_TIMEFRAMES = "getTimeframes"
	CMD_OHLCV_STREAM   = "ohlcvStream"
	CMD_GET_OHLCV      = "getOHLCV"

	// Tax / history commands
	CMD_GET_TRADE_HISTORY    = "getTradeHistory"
	CMD_GET_ORDER_HISTORY    = "getOrderHistory"
	CMD_GET_TRANSFER_HISTORY = "getTransferHistory"
)

func AllCommands() map[string]any {
	return map[string]any{
		"CMD_GET_ACCOUNT":          CMD_GET_ACCOUNT,
		"CMD_GET_BALANCES":         CMD_GET_BALANCES,
		"CMD_GET_POSITIONS":        CMD_GET_POSITIONS,
		"CMD_GET_ORDERS":           CMD_GET_ORDERS,
		"CMD_CREATE_ORDERS":        CMD_CREATE_ORDERS,
		"CMD_CANCEL_ORDERS":        CMD_CANCEL_ORDERS,
		"CMD_GET_MARKETS":          CMD_GET_MARKETS,
		"CMD_GET_MARKET":           CMD_GET_MARKET,
		"CMD_GET_TIMEFRAMES":       CMD_GET_TIMEFRAMES,
		"CMD_OHLCV_STREAM":         CMD_OHLCV_STREAM,
		"CMD_GET_OHLCV":            CMD_GET_OHLCV,
		"CMD_GET_TRADE_HISTORY":    CMD_GET_TRADE_HISTORY,
		"CMD_GET_ORDER_HISTORY":    CMD_GET_ORDER_HISTORY,
		"CMD_GET_TRANSFER_HISTORY": CMD_GET_TRANSFER_HISTORY,
	}
}
