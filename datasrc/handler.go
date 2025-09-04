package datasrc

import (
	"github.com/extism/go-pdk"
	dt "github.com/plusev-terminal/go-plugin-common/datasrc/types"
)

// DataSource interface defines the required methods for a data source plugin
type DataSource interface {
	// GetName returns the name of the data source
	GetName() string

	// GetMarkets returns all available trading markets
	GetMarkets() ([]dt.MarketMeta, error)

	// GetTimeframes returns all supported timeframes
	GetTimeframes() []dt.Timeframe

	// GetOHLCV fetches historical OHLCV data for the given parameters
	GetOHLCV(params dt.OHLCVParams) ([]dt.OHLCVRecord, error)

	// StartStream starts streaming live data (optional, return error if not supported)
	StartStream(config dt.StreamConfig) error

	// SupportsStreaming returns true if this data source supports real-time streaming
	SupportsStreaming() bool
}

// PluginHandler provides a convenient way to implement all required plugin functions
type PluginHandler struct {
	Config     DataSourceConfig
	DataSource DataSource
}

// NewPluginHandler creates a new plugin handler with the given config and data source
func NewPluginHandler(config DataSourceConfig, ds DataSource) *PluginHandler {
	return &PluginHandler{
		Config:     config,
		DataSource: ds,
	}
}

// ExportMeta implements the meta export function
func (h *PluginHandler) ExportMeta() int32 {
	return ExportMeta(h.Config)
}

// ExportGetName implements the get_name export function
func (h *PluginHandler) ExportGetName() int32 {
	return ExportName(h.DataSource.GetName())
}

// ExportListMarkets implements the list_markets export function
func (h *PluginHandler) ExportListMarkets() int32 {
	markets, err := h.DataSource.GetMarkets()
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	return ExportMarkets(markets)
}

// ExportGetTimeframes implements the get_timeframes export function
func (h *PluginHandler) ExportGetTimeframes() int32 {
	timeframes := h.DataSource.GetTimeframes()
	return ExportTimeframes(timeframes)
}

// ExportGetOHLCV implements the get_ohlcv export function
func (h *PluginHandler) ExportGetOHLCV() int32 {
	params, err := GetOHLCVParams()
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	data, err := h.DataSource.GetOHLCV(params)
	return ExportOHLCV(data, err)
}

// ExportStreamOHLCV implements the stream_ohlcv export function
func (h *PluginHandler) ExportStreamOHLCV() int32 {
	config, err := GetStreamConfig()
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	err = h.DataSource.StartStream(config)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	return 0
}
