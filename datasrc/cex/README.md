# CEX Command Parameters

This package provides type-safe parameter structs for CEX (Centralized Exchange) datasource commands.

## Architecture

**Validation happens in the terminal** before commands are forwarded to datasources (built-in or plugins). This means:

- ✅ Terminal validates all parameters using `Validate*Params()` functions
- ✅ Plugins receive guaranteed-valid parameters
- ✅ Plugins use simple `*ParamsFromMap()` extraction functions (no error handling needed)
- ✅ Cleaner plugin code with less error checking boilerplate

## Features

- **WASM-Compatible**: All helper functions use basic type assertions, no reflection required
- **Type-Safe**: Strongly typed parameters instead of `map[string]any`
- **Pre-Validated**: Plugins can assume all required fields are present and valid types
- **Simple Extraction**: Just extract, no validation needed in plugins

## Available Commands

### 1. GetMarkets

Retrieve available trading pairs from the exchange.

**Parameters**: None

### 2. GetTimeframes

Get supported timeframes/intervals for the exchange.

**Parameters**: None

### 3. OHLCVStream

Subscribe to real-time OHLCV (candlestick) data stream.

**Parameters**:
- `symbol` (string, required): Trading pair symbol
- `interval` (string, required): Timeframe/interval

**Example Usage in Plugin**:
```go
func handleOHLCVStream(params map[string]any) dt.Response {
    // No error checking needed - params are pre-validated
    p := cex.OHLCVStreamParamsFromMap(params)
    
    streamID := fmt.Sprintf("%s_%s", p.Symbol, p.Interval)
    // ...
}
```

### 4. GetOHLCV

Fetch historical OHLCV data.

**Parameters**:
- `symbol` (string, required): Trading pair symbol
- `timeframe` (string, required): Timeframe/interval
- `startTime` (time.Time, optional): Start time for historical data
- `endTime` (time.Time, optional): End time for historical data
- `limit` (int, optional): Maximum number of records to return

**Time Format Support**:
- RFC3339 string: `"2024-01-01T00:00:00Z"`
- Unix milliseconds (int64): `1704067200000`
- Unix milliseconds (float64): `1704067200000.0` (JSON number)
- `time.Time` object (Go native)

**Example Usage in Plugin**:
```go
func handleGetOHLCV(params map[string]any) dt.Response {
    // No error checking needed - params are pre-validated
    p := cex.GetOHLCVParamsFromMap(params)
    
    // Optional fields are nil if not provided
    if p.StartTime != nil {
        fmt.Printf("Starting from: %s\n", p.StartTime)
    }
    
    // Use the params...
}
```

## Terminal-Side Validation

**For terminal/built-in datasource developers:**

Use the validation functions with the global validator before forwarding commands:

```go
import (
    "github.com/plusev-terminal/go-plugin-common/datasrc/cex"
    "your-terminal/global"
)

// In datasource actor or command handler
func handleCommand(cmd string, params map[string]any) {
    // Validate before forwarding using terminal's global validator
    var err error
    switch cmd {
    case cex.CEX_CMD_OHLCV_STREAM:
        err = cex.ValidateOHLCVStreamParams(global.Validator, params)
    case cex.CEX_CMD_GET_OHLCV:
        err = cex.ValidateGetOHLCVParams(global.Validator, params)
    }
    
    if err != nil {
        // Return error to caller, don't forward to datasource
        // The validator returns structured ValidationErrors
        return datasrc.ErrorResponse(err)
    }
    
    // Now forward to datasource with validated params
    // ...
}
```

**Validation uses**:
- `go-playground/validator/v10` with struct tags
- `mapstructure` for flexible map-to-struct decoding
- Terminal's existing `global.Validator` instance (reuses custom validators)
- Returns structured `validator.ValidationErrors` for consistent error handling

## Plugin-Side Usage

**For plugin developers:**

Simply extract the parameters - no validation needed!

```go
import "github.com/plusev-terminal/go-plugin-common/datasrc/cex"

func handleOHLCVStream(params map[string]any) dt.Response {
    // Just extract - validation already done by terminal
    p := cex.OHLCVStreamParamsFromMap(params)
    
    // All fields are guaranteed valid
    streamID := fmt.Sprintf("%s_%s", p.Symbol, p.Interval)
    // ...
}

func handleGetOHLCV(params map[string]any) dt.Response {
    // Just extract - validation already done by terminal
    p := cex.GetOHLCVParamsFromMap(params)
    
    // Required fields are always set
    fmt.Printf("Symbol: %s, Timeframe: %s\n", p.Symbol, p.Timeframe)
    
    // Optional fields may be nil/zero
    if p.StartTime != nil {
        fmt.Printf("Start: %s\n", p.StartTime)
    }
    // ...
}
```

## Helper Functions

### ExtractString
```go
symbol := cex.ExtractString(params, "symbol")
// Returns empty string if not found or wrong type
```

### ExtractInt
```go
limit := cex.ExtractInt(params, "limit")
// Returns 0 if not found or wrong type
// Handles int, int64, float64 (JSON numbers)
```

### ExtractTime
```go
startTime := cex.ExtractTime(params, "startTime")
// Returns nil if not found or invalid
// Handles: RFC3339 string, time.Time, int64/float64 (unix millis)
```

## WASM Compatibility

✅ **All functions work in WASM** because they:
- Use only basic type assertions (`value.(string)`, `value.(int)`, etc.)
- Avoid reflection-based libraries like `mapstructure`
- Handle JSON deserialization edge cases (float64 for numbers)
- Provide manual field-by-field extraction

This makes them suitable for use in Extism WASM plugins where reflection may be limited or unavailable.

## Migration from Old Approach

**Old way** (validation in plugin):
```go
func handleOHLCVStream(params map[string]any) dt.Response {
    symbolVal, ok := params["symbol"]
    if !ok {
        return datasrc.ErrorResponse(fmt.Errorf("symbol is required"))
    }
    symbol, ok := symbolVal.(string)
    if !ok {
        return datasrc.ErrorResponse(fmt.Errorf("symbol must be a string"))
    }
    // ... repeat for interval ...
}
```

**New way** (validation in terminal):
```go
func handleOHLCVStream(params map[string]any) dt.Response {
    // One line - terminal already validated!
    p := cex.OHLCVStreamParamsFromMap(params)
    // Just use it
}
```

