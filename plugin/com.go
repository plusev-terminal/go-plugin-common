package plugin

import (
	"time"

	"github.com/extism/go-pdk"
)

// Command represents a request to a plugin
type Command struct {
	Name   string         `json:"name"`   // e.g., "process", "ohlcvStream", "getMarkets", "getBalance"
	Params map[string]any `json:"params"` // Flexible parameters specific to each command
}

// Response represents the result of a command execution
type Response struct {
	Result          bool   `json:"result"`
	ResponseType    string `json:"responseType,omitempty"`    // e.g. "StreamMarker"
	Data            any    `json:"data,omitempty"`            // Could be direct data or a channel for streams
	Error           string `json:"error,omitempty"`           // Error message if Success is false
	CacheForSeconds *int64 `json:"cacheForSeconds,omitempty"` // Optional: cache duration in seconds (wrapper converts to time.Duration)
}

// StreamData represents a single piece of data from a stream
type StreamData struct {
	StreamID string `json:"streamId"` // Unique identifier for this stream
	Data     any    `json:"data"`     // The actual data (e.g., OHLCV candle, orderbook update)
}

// StreamSetupRequest represents the request sent to plugin for stream setup
type StreamSetupRequest struct {
	StreamID   string         `json:"streamId"`
	StreamType string         `json:"streamType"` // "ohlcv", "orderbook", "orders", "trades", etc.
	Parameters map[string]any `json:"parameters"` // Generic parameters
}

// StreamSetupResponse represents plugin's response to stream setup request
type StreamSetupResponse struct {
	Success         bool              `json:"success"`
	WebSocketURL    string            `json:"websocketUrl"`
	Headers         map[string]string `json:"headers,omitempty"`
	Subprotocol     string            `json:"subprotocol,omitempty"`
	InitialMessages []string          `json:"initialMessages"`
	StreamContext   map[string]any    `json:"streamContext,omitempty"`
	Error           string            `json:"error,omitempty"`
}

// StreamMessageRequest represents the request sent to plugin for message processing
type StreamMessageRequest struct {
	StreamID      string         `json:"streamId"`
	ConnectionID  string         `json:"connectionId"`
	Message       []byte         `json:"message"`
	MessageType   string         `json:"messageType"` // "data", "error", "close"
	StreamContext map[string]any `json:"streamContext,omitempty"`
}

// StreamMessageResponse represents plugin's response to a stream message
type StreamMessageResponse struct {
	Success     bool   `json:"success"`
	Action      string `json:"action"`             // "ignore", "data", "reconnect", "close", "send"
	DataType    string `json:"dataType,omitempty"` // "ohlcv", "orderbook", "order_fill", etc.
	Data        any    `json:"data,omitempty"`     // Generic data payload
	SendMessage string `json:"sendMessage,omitempty"`
	Error       string `json:"error,omitempty"`
}

// StreamConnectionEvent represents a connection lifecycle event
type StreamConnectionEvent struct {
	StreamID     string `json:"streamId"`
	ConnectionID string `json:"connectionId"`
	EventType    string `json:"eventType"` // "connected", "disconnected", "error"
	Error        string `json:"error,omitempty"`
}

// StreamConnectionResponse represents plugin's response to a connection event
type StreamConnectionResponse struct {
	Success bool   `json:"success"`
	Action  string `json:"action"` // "ignore", "reconnect", "close"
	Error   string `json:"error,omitempty"`
}

// ReadCommand reads a command from plugin input (used in handle_command export)
func ReadCommand() (Command, error) {
	var cmd Command
	err := pdk.InputJSON(&cmd)
	return cmd, err
}

// WriteResponse writes a response to plugin output
func WriteResponse(resp Response) int32 {
	pdk.OutputJSON(resp)
	if resp.Result {
		return 0
	}
	return 1
}

// SuccessResponse creates a successful response with data
func SuccessResponse(data any, cacheFor ...time.Duration) Response {
	if len(cacheFor) > 0 {
		seconds := int64(cacheFor[0].Seconds())
		return Response{
			Result:          true,
			Data:            data,
			CacheForSeconds: &seconds,
		}
	}

	return Response{
		Result: true,
		Data:   data,
	}
}

// SuccessTypedResponse creates a successful response with an explicit response type.
func SuccessTypedResponse(responseType string, data any, cacheFor ...time.Duration) Response {
	resp := SuccessResponse(data, cacheFor...)
	resp.ResponseType = responseType
	return resp
}

// ErrorResponse creates an error response
func ErrorResponse(err error) Response {
	return Response{
		Result: false,
		Error:  err.Error(),
	}
}

// ErrorResponseMsg creates an error response with a message
func ErrorResponseMsg(msg string) Response {
	return Response{
		Result: false,
		Error:  msg,
	}
}
