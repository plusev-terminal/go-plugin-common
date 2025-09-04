package datasrc

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

// WebSocket connection and streaming types

// WSConnection represents a WebSocket connection from the plugin side
type WSConnection struct {
	ID  string
	URL string
}

// StreamData represents real-time streaming data
type StreamData struct {
	Symbol    string      `json:"symbol"`
	Type      string      `json:"type"` // "ticker", "ohlcv", "trade", etc.
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// TickerData represents real-time ticker information
type TickerData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
	Change float64 `json:"change"`
}

// TradeData represents individual trade data
type TradeData struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Side      string  `json:"side"` // "buy" or "sell"
	Timestamp int64   `json:"timestamp"`
}

// WebSocket Host Function Request/Response Types (must match websocket_host.go)

type WSConnectRequest struct {
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers,omitempty"`
	Subprotocol string            `json:"subprotocol,omitempty"`
}

type WSConnectResponse struct {
	Success      bool   `json:"success"`
	ConnectionID string `json:"connectionId,omitempty"`
	Error        string `json:"error,omitempty"`
}

type WSSendRequest struct {
	ConnectionID string `json:"connectionId"`
	Message      string `json:"message"`
}

type WSSendResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type WSReceiveRequest struct {
	ConnectionID string `json:"connectionId"`
	TimeoutMs    int    `json:"timeoutMs,omitempty"` // 0 = no timeout, -1 = non-blocking
}

type WSReceiveResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Timeout bool   `json:"timeout,omitempty"`
}

type WSCloseRequest struct {
	ConnectionID string `json:"connectionId"`
}

type WSCloseResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// WebSocket Client Functions for Plugins

// WSConnect establishes a WebSocket connection
func WSConnect(url string, headers map[string]string, subprotocol string) (*WSConnection, error) {
	req := WSConnectRequest{
		URL:         url,
		Headers:     headers,
		Subprotocol: subprotocol,
	}

	// Allocate memory and marshal the request to JSON
	mem, err := pdk.AllocateJSON(req)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory: %w", err)
	}

	// Call ws_connect host function
	responsePtr := wsConnect(mem.Offset())
	responseMem := pdk.FindMemory(responsePtr)
	responseData := responseMem.ReadBytes()

	var resp WSConnectResponse
	if err := json.Unmarshal(responseData, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connect response: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("WebSocket connection failed: %s", resp.Error)
	}

	return &WSConnection{
		ID:  resp.ConnectionID,
		URL: url,
	}, nil
}

// WSSend sends a message over the WebSocket connection
func (ws *WSConnection) Send(message string) error {
	req := WSSendRequest{
		ConnectionID: ws.ID,
		Message:      message,
	}

	mem, err := pdk.AllocateJSON(req)
	if err != nil {
		return fmt.Errorf("failed to allocate memory: %w", err)
	}

	responsePtr := wsSend(mem.Offset())
	responseMem := pdk.FindMemory(responsePtr)
	responseData := responseMem.ReadBytes()

	var resp WSSendResponse
	if err := json.Unmarshal(responseData, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal send response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("WebSocket send failed: %s", resp.Error)
	}

	return nil
}

// WSReceive receives a message from the WebSocket connection
func (ws *WSConnection) Receive(timeoutMs int) (string, bool, error) {
	req := WSReceiveRequest{
		ConnectionID: ws.ID,
		TimeoutMs:    timeoutMs,
	}

	mem, err := pdk.AllocateJSON(req)
	if err != nil {
		return "", false, fmt.Errorf("failed to allocate memory: %w", err)
	}

	responsePtr := wsReceive(mem.Offset())
	responseMem := pdk.FindMemory(responsePtr)
	responseData := responseMem.ReadBytes()

	var resp WSReceiveResponse
	if err := json.Unmarshal(responseData, &resp); err != nil {
		return "", false, fmt.Errorf("failed to unmarshal receive response: %w", err)
	}

	if !resp.Success {
		return "", false, fmt.Errorf("WebSocket receive failed: %s", resp.Error)
	}

	return resp.Message, resp.Timeout, nil
}

// WSClose closes the WebSocket connection
func (ws *WSConnection) Close() error {
	req := WSCloseRequest{
		ConnectionID: ws.ID,
	}

	mem, err := pdk.AllocateJSON(req)
	if err != nil {
		return fmt.Errorf("failed to allocate memory: %w", err)
	}

	responsePtr := wsClose(mem.Offset())
	responseMem := pdk.FindMemory(responsePtr)
	responseData := responseMem.ReadBytes()

	var resp WSCloseResponse
	if err := json.Unmarshal(responseData, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal close response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("WebSocket close failed: %s", resp.Error)
	}

	return nil
}

// Host function imports (these are provided by the host application)

//go:wasmimport extism:host/user ws_connect
func wsConnect(ptr uint64) uint64

//go:wasmimport extism:host/user ws_send
func wsSend(ptr uint64) uint64

//go:wasmimport extism:host/user ws_receive
func wsReceive(ptr uint64) uint64

//go:wasmimport extism:host/user ws_close
func wsClose(ptr uint64) uint64

// Helper functions for streaming implementations

// ParseStreamData parses incoming WebSocket message as StreamData
func ParseStreamData(message string) (*StreamData, error) {
	var data StreamData
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stream data: %w", err)
	}
	return &data, nil
}

// CreateTickerData creates a StreamData with ticker information
func CreateTickerData(symbol string, price, volume, change float64) *StreamData {
	return &StreamData{
		Symbol:    symbol,
		Type:      "ticker",
		Timestamp: getCurrentTimestamp(),
		Data: TickerData{
			Symbol: symbol,
			Price:  price,
			Volume: volume,
			Change: change,
		},
	}
}

// CreateTradeData creates a StreamData with trade information
func CreateTradeData(symbol string, price, quantity float64, side string) *StreamData {
	return &StreamData{
		Symbol:    symbol,
		Type:      "trade",
		Timestamp: getCurrentTimestamp(),
		Data: TradeData{
			Symbol:    symbol,
			Price:     price,
			Quantity:  quantity,
			Side:      side,
			Timestamp: getCurrentTimestamp(),
		},
	}
}

// getCurrentTimestamp returns current Unix timestamp
func getCurrentTimestamp() int64 {
	// In a WASM environment, we'd typically get this from a host function
	// For now, return a placeholder
	return 1693737600 // This would be replaced with actual time from host
}
