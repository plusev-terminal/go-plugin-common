package cex

import (
	"time"

	"github.com/plusev-terminal/go-plugin-common/datasrc/utils"
)

// GetMarketsParams contains parameters for the getMarkets command
type GetMarketsParams struct {
	// No parameters required for getMarkets
}

// GetTimeframesParams contains parameters for the getTimeframes command
type GetTimeframesParams struct {
	// No parameters required for getTimeframes
}

// OHLCVStreamParams contains parameters for the ohlcvStream command
type OHLCVStreamParams struct {
	Symbol   string `json:"symbol" mapstructure:"symbol" validate:"required"`
	Interval string `json:"interval" mapstructure:"interval" validate:"required"`
}

// GetOHLCVParams contains parameters for the getOHLCV (historical data) command
type GetOHLCVParams struct {
	Symbol    string     `json:"symbol" mapstructure:"symbol" validate:"required"`
	Timeframe string     `json:"timeframe" mapstructure:"timeframe" validate:"required"`
	StartTime *time.Time `json:"startTime,omitempty" mapstructure:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty" mapstructure:"endTime"`
	Limit     int        `json:"limit,omitempty" mapstructure:"limit"`
}

// OHLCVStreamParamsFromMap extracts OHLCVStreamParams from validated map
func OHLCVStreamParamsFromMap(data map[string]any) OHLCVStreamParams {
	return OHLCVStreamParams{
		Symbol:   utils.ExtractString(data, "symbol"),
		Interval: utils.ExtractString(data, "interval"),
	}
}

// GetOHLCVParamsFromMap extracts GetOHLCVParams from validated map
func GetOHLCVParamsFromMap(data map[string]any) GetOHLCVParams {
	return GetOHLCVParams{
		Symbol:    utils.ExtractString(data, "symbol"),
		Timeframe: utils.ExtractString(data, "timeframe"),
		StartTime: utils.ExtractTime(data, "startTime"),
		EndTime:   utils.ExtractTime(data, "endTime"),
		Limit:     utils.ExtractInt(data, "limit"),
	}
}
