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

	SetCredentials(params map[string]string) error

	// GetCredentialFields returns the credential fields required for authentication
	GetCredentialFields() ([]dt.CredentialField, error)

	// GetMarkets returns all available trading markets
	GetMarkets() ([]dt.MarketMeta, error)

	// GetTimeframes returns all supported timeframes
	GetTimeframes() []dt.Timeframe

	// GetOHLCV fetches historical OHLCV data for the given parameters
	GetOHLCV(params dt.OHLCVParams) ([]dt.OHLCVRecord, error)

	// PrepareStream prepares streaming connection setup
	PrepareStream(request dt.StreamSetupRequest) (dt.StreamSetupResponse, error)

	// HandleStreamMessage processes incoming stream messages
	HandleStreamMessage(request dt.StreamMessageRequest) (dt.StreamMessageResponse, error)

	// HandleConnectionEvent handles stream connection events
	HandleConnectionEvent(event dt.StreamConnectionEvent) (dt.StreamConnectionResponse, error)

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

func (h *PluginHandler) ExportGetCredentialFields() int32 {
	fields, err := h.DataSource.GetCredentialFields()
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	return ExportCredentialFields(fields)
}

func (h *PluginHandler) ExportSetCredentials() int32 {
	params, err := GetCredentials()
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	err = h.DataSource.SetCredentials(params)
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	return 0
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
	var request dt.StreamSetupRequest
	err := pdk.InputJSON(&request)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	response, err := h.DataSource.PrepareStream(request)
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	pdk.OutputJSON(response)
	return 0
}

// ExportHandleStreamMessage implements the handle_stream_message export function
func (h *PluginHandler) ExportHandleStreamMessage() int32 {
	var request dt.StreamMessageRequest
	err := pdk.InputJSON(&request)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	response, err := h.DataSource.HandleStreamMessage(request)
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	pdk.OutputJSON(response)
	return 0
}

// ExportStreamConnectionEvent implements the stream_connection_event export function
func (h *PluginHandler) ExportStreamConnectionEvent() int32 {
	var event dt.StreamConnectionEvent
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
