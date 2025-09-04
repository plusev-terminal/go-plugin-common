# Requester Testing Package

This package provides testing utilities for PlusEV plugins that use the `requester` package. It allows testing plugins without the need for WASM host functions by providing mock implementations.

## MockRequester

The `MockRequester` implements `requester.Interface` and can either:
1. Return predefined mock responses
2. Make real HTTP requests using Go's standard `net/http` package

## Usage in Plugin Tests

### 1. Import the testing package

```go
import (
    "testing"
    requestertesting "github.com/plusev-terminal/go-plugin-common/requester/testing"
    "github.com/plusev-terminal/go-plugin-common/requester"
)
```

### 2. Create and configure mock requester

```go
func TestYourPlugin(t *testing.T) {
    // Create a mock requester
    mockReq := requestertesting.NewMockRequester()
    
    // Set up mock responses
    mockReq.SetMockResponse("/v3/public/instruments", `{
        "success": true,
        "data": {
            "rows": [
                {
                    "symbol": "BTC_USDT",
                    "baseAsset": "BTC",
                    "quoteAsset": "USDT",
                    "status": "TRADING"
                }
            ]
        }
    }`)
    
    // Use the mock in your plugin logic
    // Your plugin should accept requester.Interface in its constructor
    plugin := NewYourPlugin(mockReq, "https://api.example.com")
    result, err := plugin.SomeMethod()
    
    // Assert results
    if err != nil {
        t.Fatalf("Plugin method failed: %v", err)
    }
    
    // Verify API calls were made
    calls := mockReq.GetCalls()
    // ... verify expected calls were made
}
```

### 3. Making your plugin testable

Structure your plugin to accept the requester interface:

```go
// In your plugin package
type YourPlugin struct {
    requester requester.Interface
    baseURL   string
}

func NewYourPlugin(req requester.Interface, baseURL string) *YourPlugin {
    return &YourPlugin{
        requester: req,
        baseURL:   baseURL,
    }
}

// For production use (in main.go):
var plugin = NewYourPlugin(requester.NewDefault(), "https://api.example.com")

// For testing:
func TestYourPlugin(t *testing.T) {
    mockReq := requestertesting.NewMockRequester()
    plugin := NewYourPlugin(mockReq, "https://api.example.com")
    // ... test plugin
}
```

## Pattern Matching

The mock requester supports various URL pattern matching:

- **Exact match**: `SetMockResponse("https://api.example.com/endpoint", response)`
- **Contains match**: `SetMockResponse("/endpoint", response)` - matches any URL containing "/endpoint"
- **Wildcard match**: `SetMockResponse("https://api.example.com/*", response)` - matches URLs starting with the prefix

## Integration Testing

If no mock responses are set, `MockRequester` will make real HTTP requests:

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    mockReq := requestertesting.NewMockRequester()
    // Don't set any mocks - will make real requests
    
    plugin := NewYourPlugin(mockReq, "https://real-api.com")
    // Test against real API
}
```

## Error Testing

Mock error responses for error handling tests:

```go
func TestErrorHandling(t *testing.T) {
    mockReq := requestertesting.NewMockRequester()
    mockReq.SetMockError("/error-endpoint", fmt.Errorf("network timeout"))
    
    plugin := NewYourPlugin(mockReq, "https://api.example.com")
    _, err := plugin.MethodThatCallsErrorEndpoint()
    
    if err == nil {
        t.Error("Expected error but got none")
    }
}
```

## Best Practices

1. **Use interfaces**: Design your plugin to accept `requester.Interface`
2. **Test both success and error cases**: Use mocks for comprehensive testing
3. **Verify API calls**: Use `GetCalls()` to ensure correct API usage
4. **Reset between tests**: Call `Reset()` to clear state between test cases
5. **Prefer unit tests**: Use mocks for fast, reliable unit tests; use real requests sparingly for integration tests
