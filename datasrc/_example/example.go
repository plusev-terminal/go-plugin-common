package datasrc

import (
	"errors"
	"time"

	m "github.com/plusev-terminal/go-plugin-common/meta"
)

// ExampleDataSource is a simple example implementation of the DataSource interface
// This can be used as a template for creating real data source plugins
type ExampleDataSource struct {
	name string
}

// NewExampleDataSource creates a new example data source
func NewExampleDataSource(name string) *ExampleDataSource {
	return &ExampleDataSource{name: name}
}

// GetName returns the name of this data source
func (e *ExampleDataSource) GetName() string {
	return e.name
}

// GetMarkets returns example markets
func (e *ExampleDataSource) GetMarkets() ([]MarketMeta, error) {
	return []MarketMeta{
		{Name: "BTCUSDT", Base: "BTC", Quote: "USDT", AssetType: "spot"},
		{Name: "ETHUSDT", Base: "ETH", Quote: "USDT", AssetType: "spot"},
		{Name: "ADAUSDT", Base: "ADA", Quote: "USDT", AssetType: "spot"},
	}, nil
}

// GetTimeframes returns common timeframes
func (e *ExampleDataSource) GetTimeframes() []Timeframe {
	return CommonTimeframes[:6] // Return first 6 common timeframes
}

// GetOHLCV returns example OHLCV data
func (e *ExampleDataSource) GetOHLCV(params OHLCVParams) ([]OHLCVRecord, error) {
	// Generate some dummy data
	now := time.Now().Unix()
	return []OHLCVRecord{
		{
			Timestamp: now - 120,
			Open:      45000.0,
			High:      45500.0,
			Low:       44800.0,
			Close:     45200.0,
			Volume:    1250.5,
		},
		{
			Timestamp: now - 60,
			Open:      45200.0,
			High:      45300.0,
			Low:       45000.0,
			Close:     45100.0,
			Volume:    980.2,
		},
	}, nil
}

// StartStream returns an error since streaming is not implemented in this example
func (e *ExampleDataSource) StartStream(config StreamConfig) error {
	return errors.New("streaming not implemented in example data source")
}

// SupportsStreaming returns false for the example data source
func (e *ExampleDataSource) SupportsStreaming() bool {
	return false
}

// CreateExampleConfig creates an example configuration for a data source plugin
func CreateExampleConfig(pluginID, name, description, author string) DataSourceConfig {
	return DataSourceConfig{
		PluginID:    pluginID,
		Name:        name,
		Description: description,
		Author:      author,
		Version:     "1.0.0",
		Repository:  "https://github.com/your-org/your-plugin",
		Tags:        []string{"example", "datasource"},
		Contacts: []m.AuthorContact{
			{Kind: "email", Value: "your-email@example.com"},
		},
		NetworkTargets: []string{
			"https://api.binance.com/*",
			"https://api.exchange.coinbase.com/*",
		},
	}
}
