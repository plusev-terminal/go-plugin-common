package exchange

import (
	"fmt"
	"time"

	tt "github.com/plusev-terminal/go-plugin-common/trading"
	"github.com/plusev-terminal/go-plugin-common/utils"
)

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
		CacheForSeconds: utils.ExtractInt("cacheFor", data),
	}
	if v, ok := data["market"].(map[string]any); ok {
		_ = utils.MapToStruct(v, &params.Market)
	}
	return params
}

// AccountBalancesParams contains parameters for the accountBalances command.
// Market is required so the plugin can select the correct account context (spot/futures/etc)
// without relying on ad-hoc fields.
type AccountBalancesParams struct {
	Market tt.Market `json:"market" mapstructure:"market" validate:"required"`
}

func (p AccountBalancesParams) Validate() error {
	if p.Market.AssetType == "" {
		// Defaulting is left to the plugin, but an empty assetType is still a valid input.
		// We primarily require presence of the market object.
	}
	return nil
}

func AccountBalancesParamsFromMap(data map[string]any) AccountBalancesParams {
	params := AccountBalancesParams{}
	if v, ok := data["market"].(map[string]any); ok {
		_ = utils.MapToStruct(v, &params.Market)
	}
	return params
}
