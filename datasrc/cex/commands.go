package cex

type CEX_CMD string

const (
	CEX_CMD_GET_MARKETS    = "getMarkets"
	CEX_CMD_GET_TIMEFRAMES = "getTimeframes"
	CEX_CMD_OHLCV_STREAM   = "ohlcvStream"
	CEX_CMD_GET_OHLCV      = "getOHLCV"
)
