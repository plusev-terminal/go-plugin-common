package datasrc

import (
	"fmt"

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

	// PrepareStream prepares streaming connection setup
	PrepareStream(config dt.StreamConfig) (dt.StreamSetup, error)

	// HandleStreamMessage processes incoming stream messages
	HandleStreamMessage(message dt.StreamMessage) (dt.StreamResponse, error)

	// HandleConnectionEvent handles stream connection events
	HandleConnectionEvent(event dt.ConnectionEvent) (dt.ConnectionResponse, error)

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

// ExportStreamOHLCV implements the stream_ohlcv export function (DEPRECATED)
func (h *PluginHandler) ExportStreamOHLCV() int32 {
	// This is now deprecated - use the new callback-based system
	pdk.SetError(fmt.Errorf("stream_ohlcv is deprecated - use prepare_stream, handle_stream_message, and stream_connection_event"))
	return 1
}

// ExportPrepareStream implements the prepare_stream export function
func (h *PluginHandler) ExportPrepareStream() int32 {
	config, err := GetStreamConfig()
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	setup, err := h.DataSource.PrepareStream(config)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	pdk.OutputJSON(setup)
	return 0
}

// ExportHandleStreamMessage implements the handle_stream_message export function
func (h *PluginHandler) ExportHandleStreamMessage() int32 {
	var message dt.StreamMessage
	err := pdk.InputJSON(&message)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	response, err := h.DataSource.HandleStreamMessage(message)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	pdk.OutputJSON(response)
	return 0
}

// ExportStreamConnectionEvent implements the stream_connection_event export function
func (h *PluginHandler) ExportStreamConnectionEvent() int32 {
	var event dt.ConnectionEvent
	err := pdk.InputJSON(&event)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	response, err := h.DataSource.HandleConnectionEvent(event)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	pdk.OutputJSON(response)
	return 0
}
