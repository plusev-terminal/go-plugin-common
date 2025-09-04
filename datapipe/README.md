# DataPipe Plugin Package

This package provides a standardized interface for building DataPipe plugin nodes for the PlusEV Terminal.

## Overview

DataPipe plugins are WASM-based nodes that process data in visual pipelines. They receive input data, process it according to their configuration, and output results. This package handles all the WASM lifecycle and exports automatically - plugin developers just implement the `DataPipePlugin` interface.

## Quick Start

```go
package main

import (
	"github.com/plusev-terminal/go-plugin-common/datapipe"
	dt "github.com/plusev-terminal/go-plugin-common/datapipe/types"
	m "github.com/plusev-terminal/go-plugin-common/meta"
	"github.com/plusev-terminal/go-plugin-common/plugin"
)

type MyStrategyPlugin struct {
	config *plugin.ConfigStore
}

func (p *MyStrategyPlugin) GetMeta() m.Meta {
	return m.Meta{
		PluginID:    "my-strategy",
		Name:        "My Trading Strategy",
		AppID:       "datapipes",
		Category:    "Strategies",
		Description: "My awesome trading strategy",
		Author:      "Your Name",
		Version:     "1.0.0",
		Resources: m.ResourceAccess{
			// Define resource requirements
		},
	}
}

func (p *MyStrategyPlugin) GetConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Label:       "Fast MA Period",
			Name:        "fast_period",
			Type:        "number",
			Required:    true,
			Default:     10,
			Description: "Fast moving average period",
		},
	}
}

func (p *MyStrategyPlugin) OnInit(config *plugin.ConfigStore) error {
	p.config = config
	// Initialize and validate configuration
	return nil
}

func (p *MyStrategyPlugin) OnShutdown() error {
	// Cleanup resources
	return nil
}

func (p *MyStrategyPlugin) GetRateLimits() []plugin.RateLimit {
	// Define rate limits if needed
	return nil
}

func (p *MyStrategyPlugin) RegisterCommands(router *plugin.CommandRouter) {
	router.Register(datapipe.CMD_PROCESS, p.handleProcess)
	router.Register(datapipe.CMD_GET_CONNECTIONS, p.handleGetConnections)
	router.Register(datapipe.CMD_GET_NODE_META, p.handleGetNodeMeta)
}

// Command Handlers

func (p *MyStrategyPlugin) handleGetConnections(_ map[string]any) plugin.Response {
	return plugin.SuccessResponse(dt.Connections{
		Inputs: []dt.NodePort{
			{Name: "ohlcv", DataTypes: []dt.DataType{dt.DataTypeOHLCVRecord}},
		},
		Outputs: []dt.NodePort{
			{Name: "signal", DataTypes: []dt.DataType{dt.DataTypeSignal}},
		},
	})
}

func (p *MyStrategyPlugin) handleGetGuiDefinition(_ map[string]any) plugin.Response {
	// Optional: Return custom UI controls, or return nil data for default rendering
	return plugin.SuccessResponse(dt.GuiDefinition{
		Controls: []dt.GuiControl{
			{
				Label: "Fast MA Period",
				Name:  "fast_period",
				Type:  "number",
				Options: map[string]any{
					"min": 1,
					"max": 100,
				},
			},
		},
	})
}

func (p *MyStrategyPlugin) handleProcess(params map[string]any) plugin.Response {
	// Get configuration
	fastPeriod := p.config.GetNumber("fast_period")
	
	// Get input data from params (keyed by port name)
	inputs, ok := params["inputs"].(map[string]any)
	if !ok {
		return plugin.ErrorResponseMsg("missing inputs parameter")
	}

	// Get input data from the "ohlcv" port
	ohlcvRawData, ok := inputs["ohlcv"]
	if !ok {
		return plugin.ErrorResponseMsg("missing ohlcv input data")
	}
	
	// Process data
	// ... your logic here ...
	
	// Return output (keyed by output port name)
	return plugin.SuccessResponse(map[string]any{
		"signal": map[string]any{
			"type":     "buy",
			"strength": 0.8,
		},
	})
}

func init() {
	// Register plugin - MUST be in init(), not main()
	plugin.RegisterPlugin(&MyStrategyPlugin{})
}

func main() {
	// Required for WASM, but can be empty
}
```

## Building

Build your plugin for WASM:

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm main.go
```

## Interface Methods

### `GetMeta() DataPipeMeta`

Returns plugin metadata including:
- Basic info (ID, name, author, version, etc.)
- Input ports with their accepted data types
- Output ports with their output data types
- Resource requirements (network, filesystem, etc.)

### `GetConfigFields() []ConfigField`

Returns configuration fields for the plugin. These define what configuration the plugin needs.

### `RegisterCommands(router *CommandRouter)`

Registers command handlers for the plugin. Standard commands include:
- `CMD_PROCESS`: Main processing function for data flow
- `CMD_GET_CONNECTIONS`: Returns input/output port definitions
- `CMD_GET_NODE_META`: Returns custom UI controls, inputs and outputs

### `OnInit(config *ConfigStore) error`

Called when the plugin is initialized. Load and validate configuration here.

### `OnShutdown() error`

Called when the plugin is being shut down. Clean up resources.

### `Process(req ProcessRequest) ProcessResponse`

Main processing function. Called for each data flow through the node.

- `req.Input`: Map of input data keyed by input port name
- `req.Config`: User's configuration values for this node instance
- Returns: `ProcessResponse` with output data keyed by output port name

## Data Types

### Standard Data Types

- `DataTypeOHLCVRecord`: Candlestick/OHLCV data
- `DataTypeSignal`: Trading signals
- `DataTypeStartSignal`: Pipeline start trigger

### Helper Types

- `OHLCVRecord`: Struct with `ToMap()`/`FromMap()` helpers
- `Signal`: Struct with `ToMap()`/`FromMap()` helpers
- `StartSignal`: Struct with `ToMap()`/`FromMap()` helpers

## Configuration Store

The `ConfigStore` provides typed accessors for configuration values:

```go
// String values
apiKey := config.GetString("api_key")
apiKey := config.GetStringOr("api_key", "default")

// Numeric values
period := config.GetNumber("period")
period := config.GetNumberOr("period", 14)

// Boolean values
enabled := config.GetBool("enabled")

// Raw values
value := config.Get("custom_field")

// Check existence
if config.Has("optional_field") {
	// ...
}
```

## Helper Functions

### Response Helpers

```go
// Success response
return datapipe.SuccessResponse(map[string]any{
	"output_port": data,
})

// Error responses
return datapipe.ErrorResponse(err)
return datapipe.ErrorResponseMsg("Something went wrong")
```

## Best Practices

1. **Always use `init()`**: Call `RegisterPlugin()` in `init()`, never in `main()`
2. **Validate configuration**: Check configuration in `OnInit()` and return errors for invalid values
3. **Handle errors gracefully**: Return error responses instead of panicking
4. **Clean up resources**: Implement `OnShutdown()` to clean up connections, goroutines, etc.
5. **Use typed accessors**: Use `GetNumber()`, `GetString()`, etc. instead of raw `Get()`
6. **Document your ports**: Use clear names and specify supported data types

## Comparison with DataSource Plugins

DataPipe plugins are similar to DataSource plugins but have some key differences:

| Feature | DataSource | DataPipe |
|---------|-----------|----------|
| Purpose | Fetch external data | Process data in pipelines |
| Interface | Command-based | Direct processing |
| Inputs | Command parameters | Data from connected nodes |
| Outputs | Command responses | Data to connected nodes |
| State | Per-connection | Per-node instance |
| Lifecycle | Command execution | Data flow events |

DataPipe plugins process data **synchronously** as it flows through the pipeline, while DataSource plugins respond to **commands** (like `getOHLCV`, `getMarkets`, etc.).

## See Also

- DataSource plugin package: `github.com/plusev-terminal/go-plugin-common/datasrc`
- Meta package: `github.com/plusev-terminal/go-plugin-common/meta`
- Extism PDK: `github.com/extism/go-pdk`
