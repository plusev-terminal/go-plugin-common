package cex

import (
	"time"

	"github.com/plusev-terminal/go-plugin-common/utils"
)

// OHLCVStreamParams contains parameters for the ohlcvStream command
type OHLCVStreamParams struct {
	Symbol    string `json:"symbol" mapstructure:"symbol" validate:"required"`
	Timeframe string `json:"timeframe" mapstructure:"timeframe" validate:"required"`
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
		Symbol:    utils.GetValue[string]("symbol", data),
		Timeframe: utils.GetValue[string]("timeframe", data),
	}
}

// GetOHLCVParamsFromMap extracts GetOHLCVParams from validated map
func GetOHLCVParamsFromMap(data map[string]any) GetOHLCVParams {
	return GetOHLCVParams{
		Symbol:    utils.GetValue[string]("symbol", data),
		Timeframe: utils.GetValue[string]("timeframe", data),
		StartTime: utils.ExtractTime(data, "startTime"),
		EndTime:   utils.ExtractTime(data, "endTime"),
		Limit:     utils.ExtractInt(data, "limit"),
	}
}
