package plugin

import (
	"github.com/extism/go-pdk"
)

// StreamHandler is the interface that plugin developers implement to handle WebSocket streaming
// This is separate from the main DataSourcePlugin interface because not all data sources need streaming
type StreamHandler interface {
	// HandleStreamMessage processes incoming WebSocket messages
	// Return action="data" to push data to consumers, or action="ignore" to skip
	HandleStreamMessage(request StreamMessageRequest) (StreamMessageResponse, error)

	// HandleConnectionEvent handles WebSocket connection lifecycle events
	// Return action="reconnect" to request reconnection, or action="ignore" to do nothing
	HandleConnectionEvent(event StreamConnectionEvent) (StreamConnectionResponse, error)
}

// Global stream handler registered by RegisterStreamHandler
var registeredStreamHandler StreamHandler

// RegisterStreamHandler registers a StreamHandler implementation and enables stream WASM exports
// Call this in init() after RegisterPlugin if your plugin supports WebSocket streaming
//
// Example:
//
//	type MyPlugin struct {
//	    client *MyClient // implements StreamHandler
//	}
//
//	func init() {
//	    plugin := &MyPlugin{}
//	    datasrc.RegisterPlugin(plugin)
//	    datasrc.RegisterStreamHandler(plugin.client)
//	}
//
// After calling this, the plugin will expose handle_stream_message and handle_connection_event
// WASM exports that the host will call to deliver WebSocket messages and connection events.
func RegisterStreamHandler(handler StreamHandler) {
	registeredStreamHandler = handler
}

// ============================================================================
// WASM Exports for Stream Handling - Auto-generated
// ============================================================================

//go:wasmexport handle_stream_message
func handle_stream_message() int32 {
	// Check if stream handler is registered
	if registeredStreamHandler == nil {
		pdk.OutputJSON(StreamMessageResponse{
			Success: false,
			Action:  "ignore",
			Error:   "stream handler not registered",
		})
		return 1
	}

	// Read the incoming request
	var req StreamMessageRequest
	if err := pdk.InputJSON(&req); err != nil {
		pdk.OutputJSON(StreamMessageResponse{
			Success: false,
			Action:  "ignore",
			Error:   "failed to parse stream message request",
		})
		return 1
	}

	// Call the registered handler
	resp, err := registeredStreamHandler.HandleStreamMessage(req)
	if err != nil {
		pdk.OutputJSON(StreamMessageResponse{
			Success: false,
			Action:  "ignore",
			Error:   err.Error(),
		})
		return 1
	}

	// Write the response
	pdk.OutputJSON(resp)
	return 0
}

//go:wasmexport handle_connection_event
func handle_connection_event() int32 {
	// Check if stream handler is registered
	if registeredStreamHandler == nil {
		pdk.OutputJSON(StreamConnectionResponse{
			Success: false,
			Action:  "ignore",
			Error:   "stream handler not registered",
		})
		return 1
	}

	// Read the incoming event
	var event StreamConnectionEvent
	if err := pdk.InputJSON(&event); err != nil {
		pdk.OutputJSON(StreamConnectionResponse{
			Success: false,
			Action:  "ignore",
			Error:   "failed to parse connection event",
		})
		return 1
	}

	// Call the registered handler
	resp, err := registeredStreamHandler.HandleConnectionEvent(event)
	if err != nil {
		pdk.OutputJSON(StreamConnectionResponse{
			Success: false,
			Action:  "ignore",
			Error:   err.Error(),
		})
		return 1
	}

	// Write the response
	pdk.OutputJSON(resp)
	return 0
}

// ============================================================================
// Helper Functions
// ============================================================================

// DefaultConnectionEventHandler provides standard reconnection logic for most plugins
// Use this as a reference or call it directly from your HandleConnectionEvent implementation
//
// Example:
//
//	func (c *Client) HandleConnectionEvent(event StreamConnectionEvent) (StreamConnectionResponse, error) {
//	    c.log.InfoWithData("Connection event", map[string]any{"type": event.EventType})
//	    return datasrc.DefaultConnectionEventHandler(event), nil
//	}
func DefaultConnectionEventHandler(event StreamConnectionEvent) StreamConnectionResponse {
	switch event.EventType {
	case "connected", "connecting":
		// Connection established or in progress - no action needed
		return StreamConnectionResponse{
			Success: true,
			Action:  "ignore",
		}
	case "disconnected", "error":
		// Connection lost or error - request reconnection
		return StreamConnectionResponse{
			Success: true,
			Action:  "reconnect",
		}
	default:
		// Unknown event type - ignore
		return StreamConnectionResponse{
			Success: true,
			Action:  "ignore",
		}
	}
}

// StreamResponse is a helper to create successful data responses
func StreamResponse(dataType string, data any) StreamMessageResponse {
	return StreamMessageResponse{
		Success:  true,
		Action:   "data",
		DataType: dataType,
		Data:     data,
	}
}

// IgnoreResponse is a helper to create ignore responses (for messages that don't need processing)
func IgnoreResponse() StreamMessageResponse {
	return StreamMessageResponse{
		Success: true,
		Action:  "ignore",
	}
}

// SendResponse is a helper to create send responses (to send a message back to the WebSocket)
func SendResponse(message string) StreamMessageResponse {
	return StreamMessageResponse{
		Success:     true,
		Action:      "send",
		SendMessage: message,
	}
}

// ReconnectResponse is a helper to request reconnection
func ReconnectResponse(reason string) StreamMessageResponse {
	return StreamMessageResponse{
		Success: true,
		Action:  "reconnect",
		Error:   reason,
	}
}
