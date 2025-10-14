# Plugin Caching Guide

This guide explains how plugins can leverage the datasrc caching system to improve performance.

## Overview

Plugins can instruct the DataSource Manager to cache their responses by setting the `CacheForSeconds` field in the response. The wrapper automatically converts this to a `time.Duration` and the manager handles all caching operations.

## How to Use Caching in Plugins

### Setting Cache Duration

Simply set the `cacheForSeconds` field in your response JSON:

**Go Plugin Example:**
```go
package main

import (
    "encoding/json"
    dt "github.com/plusev-terminal/go-plugin-common/datasrc/types"
)

func handleCommand(cmdJSON []byte) []byte {
    var cmd dt.Command
    json.Unmarshal(cmdJSON, &cmd)
    
    switch cmd.Name {
    case "getMarkets":
        markets := fetchMarkets()
        
        // Cache for 1 hour (3600 seconds)
        cacheSeconds := int64(3600)
        
        resp := dt.Response{
            Result:          true,
            Data:            markets,
            CacheForSeconds: &cacheSeconds,
        }
        
        respJSON, _ := json.Marshal(resp)
        return respJSON
        
    case "getTimeframes":
        timeframes := getTimeframes()
        
        // Cache for 24 hours (86400 seconds)
        cacheSeconds := int64(86400)
        
        resp := dt.Response{
            Result:          true,
            Data:            timeframes,
            CacheForSeconds: &cacheSeconds,
        }
        
        respJSON, _ := json.Marshal(resp)
        return respJSON
        
    case "getOHLCV":
        // Don't cache - omit CacheForSeconds
        ohlcv := fetchOHLCV(cmd.Params)
        
        resp := dt.Response{
            Result: true,
            Data:   ohlcv,
            // No CacheForSeconds = no caching
        }
        
        respJSON, _ := json.Marshal(resp)
        return respJSON
    }
}
```

**Rust Plugin Example:**
```rust
use serde::{Deserialize, Serialize};

#[derive(Serialize)]
struct Response {
    result: bool,
    data: serde_json::Value,
    #[serde(skip_serializing_if = "Option::is_none")]
    cache_for_seconds: Option<i64>,
}

fn handle_command(cmd_json: &[u8]) -> Vec<u8> {
    let cmd: Command = serde_json::from_slice(cmd_json).unwrap();
    
    match cmd.name.as_str() {
        "getMarkets" => {
            let markets = fetch_markets();
            
            let resp = Response {
                result: true,
                data: serde_json::to_value(markets).unwrap(),
                cache_for_seconds: Some(3600), // 1 hour
            };
            
            serde_json::to_vec(&resp).unwrap()
        }
        
        "getTimeframes" => {
            let timeframes = get_timeframes();
            
            let resp = Response {
                result: true,
                data: serde_json::to_value(timeframes).unwrap(),
                cache_for_seconds: Some(86400), // 24 hours
            };
            
            serde_json::to_vec(&resp).unwrap()
        }
        
        "getOHLCV" => {
            let ohlcv = fetch_ohlcv(&cmd.params);
            
            let resp = Response {
                result: true,
                data: serde_json::to_value(ohlcv).unwrap(),
                cache_for_seconds: None, // No caching
            };
            
            serde_json::to_vec(&resp).unwrap()
        }
        
        _ => panic!("Unknown command")
    }
}
```

**JavaScript/TypeScript Plugin Example:**
```typescript
interface Response {
  result: boolean;
  data?: any;
  error?: string;
  cacheForSeconds?: number;
}

function handleCommand(cmdJSON: Uint8Array): Uint8Array {
  const cmd = JSON.parse(new TextDecoder().decode(cmdJSON));
  
  switch (cmd.name) {
    case 'getMarkets': {
      const markets = fetchMarkets();
      
      const resp: Response = {
        result: true,
        data: markets,
        cacheForSeconds: 3600, // 1 hour
      };
      
      return new TextEncoder().encode(JSON.stringify(resp));
    }
    
    case 'getTimeframes': {
      const timeframes = getTimeframes();
      
      const resp: Response = {
        result: true,
        data: timeframes,
        cacheForSeconds: 86400, // 24 hours
      };
      
      return new TextEncoder().encode(JSON.stringify(resp));
    }
    
    case 'getOHLCV': {
      const ohlcv = fetchOHLCV(cmd.params);
      
      const resp: Response = {
        result: true,
        data: ohlcv,
        // No cacheForSeconds = no caching
      };
      
      return new TextEncoder().encode(JSON.stringify(resp));
    }
  }
}
```

## Cache Duration Guidelines

### Recommended Cache Durations

| Data Type | Recommended Duration | Seconds | Reasoning |
|-----------|---------------------|---------|-----------|
| Markets List | 1-6 hours | 3600-21600 | Markets change infrequently |
| Timeframes | 24 hours - 7 days | 86400-604800 | Almost never change |
| Exchange Info | 6-24 hours | 21600-86400 | Fees/limits change rarely |
| Historical OHLCV | 1 hour | 3600 | Completed candles don't change |
| Account Balance | Don't cache | - | Changes frequently |
| Live OHLCV | Don't cache | - | Real-time data |
| Order Placement | Don't cache | - | Write operation |

### When NOT to Cache

- ❌ **Streams** - Real-time data should never be cached
- ❌ **Write operations** - Place order, cancel order, etc.
- ❌ **User-specific real-time data** - Balance, positions, open orders
- ❌ **Recent incomplete candles** - Last 1-2 OHLCV candles might change

### When TO Cache

- ✅ **Metadata** - Markets, timeframes, exchange info
- ✅ **Historical data** - Completed OHLCV candles
- ✅ **Static configuration** - Fee schedules, trading rules
- ✅ **Lookup tables** - Symbol mappings, asset info

## How the System Works

### 1. Plugin Returns Response with Cache Duration

```json
{
  "result": true,
  "data": [...],
  "cacheForSeconds": 3600
}
```

### 2. Wrapper Converts to time.Duration

The plugin wrapper automatically converts `cacheForSeconds` to Go's `time.Duration`:

```go
if response.CacheForSeconds != nil && *response.CacheForSeconds > 0 {
    cacheDuration := time.Duration(*response.CacheForSeconds) * time.Second
    response.CacheFor = &cacheDuration
}
```

### 3. Manager Caches Response

The DataSource Manager:
1. Generates a unique cache key: `datasrc:conn_{id}:cmd_{name}:{params}`
2. Stores the entire response in bbolt database
3. Returns cached responses on subsequent calls
4. Automatically expires entries after TTL

### 4. Consumer Gets Cached Data

Consumers don't need to know about caching - it's completely transparent:

```go
// First call - cache miss
resp1 := datasrc.Request(connID, cmd) // Executes plugin

// Second call - cache hit!
resp2 := datasrc.Request(connID, cmd) // Returns cached data
```

## Cache Invalidation

### Force Refresh

Consumers can bypass cache with `forceRefresh`:

```go
cmd := types.Command{
    Name:         "getMarkets",
    ForceRefresh: true, // Bypass cache
}
resp := datasrc.Request(connID, cmd)
```

This is useful when:
- User explicitly requests fresh data
- Debugging stale cache issues
- After configuration changes

### Manager-Level Invalidation

The manager provides methods to clear cache:

```go
// Clear all cache for a connection
manager.ClearConnectionCache(connectionID)

// Clear cache for specific command
manager.ClearCommandCache(connectionID, "getMarkets")
```

## Best Practices

### 1. Use Appropriate TTLs

Don't over-cache or under-cache:

```go
// ✅ Good - Markets cached for 1 hour
cacheSeconds := int64(3600)

// ❌ Bad - Markets cached for 1 second (defeats purpose)
cacheSeconds := int64(1)

// ❌ Bad - Markets cached for 30 days (probably stale)
cacheSeconds := int64(2592000)
```

### 2. Don't Cache Errors

Only set `cacheForSeconds` on successful responses:

```go
if err != nil {
    return dt.Response{
        Result: false,
        Error:  err.Error(),
        // Don't set CacheForSeconds
    }
}

return dt.Response{
    Result:          true,
    Data:            data,
    CacheForSeconds: &cacheSeconds, // Only cache success
}
```

### 3. Consider Parameter Variations

Different parameters create different cache keys:

```go
// These are cached separately:
// datasrc:conn_1:cmd_getOHLCV:symbol=BTCUSDT:interval=1m
// datasrc:conn_1:cmd_getOHLCV:symbol=BTCUSDT:interval=5m
```

### 4. Document Your Caching Strategy

Let users know what's cached:

```go
// GetMarkets returns all available markets.
// Results are cached for 1 hour to reduce API load.
func GetMarkets() []Market { ... }
```

## Testing

### Test Cache Behavior

```go
// Test that caching works
resp1 := plugin.HandleCommand(cmd)
resp2 := plugin.HandleCommand(cmd)

// Check CacheForSeconds is set
assert.NotNil(resp1.CacheForSeconds)
assert.Equal(int64(3600), *resp1.CacheForSeconds)
```

### Test Cache Bypass

```go
// Test force refresh
cmd.ForceRefresh = true
resp := plugin.HandleCommand(cmd)
// Should execute fresh even if cached
```

## Common Time Durations

For convenience, here are common durations in seconds:

```go
const (
    OneMinute   = 60
    FiveMinutes = 300
    TenMinutes  = 600
    ThirtyMins  = 1800
    OneHour     = 3600
    SixHours    = 21600
    TwelveHours = 43200
    OneDay      = 86400
    OneWeek     = 604800
)

// Usage:
cacheSeconds := int64(OneHour)
```

Or use time calculations:

```go
import "time"

// 1 hour
cacheSeconds := int64(time.Hour.Seconds())

// 24 hours
cacheSeconds := int64((24 * time.Hour).Seconds())
```

## Monitoring

While the system doesn't expose built-in cache metrics yet, you can:

1. Check cache file size: `ls -lh datasource_cache.db`
2. Monitor API call reduction
3. Track response times (cached should be faster)

## Future Enhancements

Potential features being considered:

- Cache statistics API
- Per-user cache isolation
- Cache warming strategies
- Conditional caching based on data freshness
- Cache size limits and eviction policies

## Summary

Caching in plugins is simple:

1. ✅ Set `cacheForSeconds` in your response
2. ✅ Use appropriate TTLs for your data
3. ✅ Don't cache errors or real-time data
4. ✅ Test your caching behavior

The system handles everything else automatically!
