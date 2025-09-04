# DataSource Plugin Common Package

This package provides common types, interfaces, and utilities for creating DataSource plugins for PlusEV Terminal.

## Overview

The DataSource plugin system allows you to create plugins that provide market data from various exchanges or data providers. This common package reduces boilerplate and ensures consistency across plugins.

## Quick Start

### 1. Create Your Data Source Implementation

```go
package main

import (
    "github.com/plusev-terminal/go-plugin-common/datasrc"
    m "github.com/plusev-terminal/go-plugin-common/meta"
)

// MyExchange implements the datasrc.DataSource interface
type MyExchange struct {
    name string
}

func (e *MyExchange) GetName() string {
    return e.name
}

func (e *MyExchange) GetMarkets() ([]datasrc.MarketMeta, error) {
    // Fetch markets from your exchange API
    return []datasrc.MarketMeta{
        {Name: "BTCUSDT", Base: "BTC", Quote: "USDT", AssetType: "spot"},
        // ... more markets
    }, nil
}

func (e *MyExchange) GetTimeframes() []datasrc.Timeframe {
    // Return supported timeframes
    return datasrc.CommonTimeframes[:6] // Use common ones or define your own
}

func (e *MyExchange) GetOHLCV(params datasrc.OHLCVParams) ([]datasrc.OHLCVRecord, error) {
    // Fetch OHLCV data from your exchange API
    // Use params.Symbol, params.Timeframe, params.StartTime, etc.
    return []datasrc.OHLCVRecord{
        // ... your OHLCV data
    }, nil
}

func (e *MyExchange) StartStream(config datasrc.StreamConfig) error {
    // Implement streaming if supported, or return error
    return errors.New("streaming not implemented")
}
```

### 2. Create Plugin Configuration

```go
func main() {
    // Configure your plugin
    config := datasrc.DataSourceConfig{
        PluginID:    "my-exchange-plugin",
        Name:        "My Exchange Data Source",
        Description: "Provides market data from My Exchange",
        Author:      "Your Name",
        Version:     "1.0.0",
        Repository:  "https://github.com/yourorg/my-exchange-plugin",
        Tags:        []string{"exchange", "crypto", "data"},
        Contacts: []m.AuthorContact{
            {Kind: "email", Value: "you@example.com"},
        },
        NetworkTargets: []string{
            "https://api.myexchange.com/*",
        },
    }

    // Create data source instance
    dataSource := &MyExchange{name: "MyExchange"}
    
    // Create plugin handler
    handler := datasrc.NewPluginHandler(config, dataSource)
    
    // The handler provides all the export functions you need
    _ = handler // Handler is used in export functions below
}
```

### 3. Export Required Functions

```go
//go:wasmexport meta
func meta() int32 {
    return handler.ExportMeta()
}

//go:wasmexport get_name
func getName() int32 {
    return handler.ExportGetName()
}

//go:wasmexport list_markets
func listMarkets() int32 {
    return handler.ExportListMarkets()
}

//go:wasmexport get_timeframes
func getTimeframes() int32 {
    return handler.ExportGetTimeframes()
}

//go:wasmexport get_ohlcv
func getOHLCV() int32 {
    return handler.ExportGetOHLCV()
}

//go:wasmexport stream_ohlcv
func streamOHLCV() int32 {
    return handler.ExportStreamOHLCV()
}
```

That's it! Your plugin is ready to be built with TinyGo.

## Types

### MarketMeta
Represents a trading market/pair:
- `Name`: Trading pair name (e.g., "BTCUSDT")
- `Base`: Base asset (e.g., "BTC")
- `Quote`: Quote asset (e.g., "USDT")
- `AssetType`: Asset type (e.g., "spot", "futures")

### Timeframe
Represents a supported timeframe:
- `Label`: Human-readable label (e.g., "1m", "5m")
- `ApiValue`: Value used for API calls
- `Interval`: Interval in seconds

### OHLCVRecord
Represents OHLCV (candlestick) data:
- `Timestamp`: Unix timestamp
- `Open`, `High`, `Low`, `Close`: Price data
- `Volume`: Trading volume

## Network Access

Specify the network targets your plugin needs access to in the `NetworkTargets` field of your config. The plugin system will only allow requests to these patterns.

Example:
```go
NetworkTargets: []string{
    "https://api.binance.com/*",
    "https://api.exchange.coinbase.com/*",
},
```

## Common Timeframes

The package provides `CommonTimeframes` which includes standard timeframes that most exchanges support:
- 1m, 5m, 15m, 30m, 1h, 4h, 1d

You can use these directly or define your own timeframes.

## Error Handling

The handler automatically handles error reporting to the host application. Just return errors from your DataSource methods and they will be properly propagated.

## Building

Build your plugin with TinyGo:
```bash
tinygo build -o your-plugin.wasm -target wasip1 -buildmode=c-shared .
```

Make sure to use the correct build flags for compatibility with the PlusEV plugin system.
