package datasrc

import (
	"fmt"

	"github.com/extism/go-pdk"
	dt "github.com/plusev-terminal/go-plugin-common/datasrc/types"
	m "github.com/plusev-terminal/go-plugin-common/meta"
)

// DataSourceConfig contains the configuration for a data source plugin
type DataSourceConfig struct {
	PluginID    string
	Name        string
	Description string
	Author      string
	Version     string
	Repository  string
	Tags        []string
	Contacts    []m.AuthorContact
	// Network targets that this plugin needs access to
	NetworkTargets []string
}

// CreateMeta creates a properly formatted Meta struct for data source plugins
func CreateMeta(config DataSourceConfig) m.Meta {
	networkRules := make([]m.NetworkTargetRule, len(config.NetworkTargets))
	for i, target := range config.NetworkTargets {
		networkRules[i] = m.NetworkTargetRule{Pattern: target}
	}

	return m.Meta{
		PluginID:    config.PluginID,
		Name:        config.Name,
		AppID:       "datasrc",
		Category:    "DataSource",
		Description: config.Description,
		Author:      config.Author,
		Version:     config.Version,
		Repository:  config.Repository,
		Tags:        config.Tags,
		Contacts:    config.Contacts,
		Resources: m.ResourceAccess{
			AllowedNetworkTargets: networkRules,
			FsWriteAccess:         nil,
			StdoutAccess:          true,
			StderrAccess:          true,
		},
	}
}

// ExportMeta exports the meta function for a data source plugin
func ExportMeta(config DataSourceConfig) int32 {
	meta := CreateMeta(config)
	pdk.OutputJSON(meta)
	return 0
}

// ExportName exports the get_name function with the given name
func ExportName(name string) int32 {
	pdk.OutputString(name)
	return 0
}

// ExportMarkets exports the list_markets function with the given markets
func ExportMarkets(markets []dt.MarketMeta) int32 {
	pdk.OutputJSON(markets)
	return 0
}

// ExportTimeframes exports the get_timeframes function with the given timeframes
func ExportTimeframes(timeframes []dt.Timeframe) int32 {
	pdk.OutputJSON(timeframes)
	return 0
}

// GetOHLCVParams reads and parses OHLCV parameters from plugin input
func GetOHLCVParams() (dt.OHLCVParams, error) {
	var params dt.OHLCVParams
	err := pdk.InputJSON(&params)
	return params, err
}

// ExportOHLCV exports OHLCV data, handling errors appropriately
func ExportOHLCV(data []dt.OHLCVRecord, err error) int32 {
	if err != nil {
		pdk.SetError(fmt.Errorf("failed to get OHLCV data: %w", err))
		return 1
	}

	pdk.OutputJSON(data)
	return 0
}

// GetStreamConfig reads and parses stream configuration from plugin input
func GetStreamConfig() (dt.StreamConfig, error) {
	var config dt.StreamConfig
	err := pdk.InputJSON(&config)
	return config, err
}

// Common timeframes that most exchanges support
var CommonTimeframes = []dt.Timeframe{
	{Label: "1m", ApiValue: "1m", Interval: 60},
	{Label: "5m", ApiValue: "5m", Interval: 300},
	{Label: "15m", ApiValue: "15m", Interval: 900},
	{Label: "30m", ApiValue: "30m", Interval: 1800},
	{Label: "1h", ApiValue: "1h", Interval: 3600},
	{Label: "4h", ApiValue: "4h", Interval: 14400},
	{Label: "1d", ApiValue: "1d", Interval: 86400},
}
