package stream

import (
	"fmt"
	"strings"
)

// StreamMarker is returned as Response.Data when a command wants the host
// to establish a WebSocket connection and manage streaming.
//
// Strict contract: the host expects this typed JSON shape.
type StreamMarker struct {
	Stream bool `json:"_stream"`

	StreamID     string            `json:"streamID"`
	WebSocketURL string            `json:"websocketUrl"`
	Headers      map[string]string `json:"headers,omitempty"`
	Subprotocol  string            `json:"subprotocol,omitempty"`

	InitialMessages []string `json:"initialMessages,omitempty"`

	// StreamContext is persisted by the host per stream and forwarded back
	// to the plugin on every handle_stream_message callback.
	StreamContext map[string]any `json:"streamContext,omitempty"`

	// Heartbeat describes how the host should handle keepalive for this stream.
	Heartbeat *StreamHeartbeatSpec `json:"heartbeat,omitempty"`
}

func (m StreamMarker) Validate() error {
	if !m.Stream {
		return fmt.Errorf("_stream must be true")
	}
	if strings.TrimSpace(m.StreamID) == "" {
		return fmt.Errorf("streamID is required")
	}
	if strings.TrimSpace(m.WebSocketURL) == "" {
		return fmt.Errorf("websocketUrl is required")
	}
	return nil
}

// StreamHeartbeatSpec describes keepalive behavior for stream connections.
//
// Note: exchanges vary: some use WS control-frame ping/pong (transport),
// others define app-level ping/pong JSON messages. This allows plugins to
// declare the protocol details while the host owns the mechanics.
type StreamHeartbeatSpec struct {
	App       *AppHeartbeatSpec       `json:"app,omitempty"`
	Transport *TransportHeartbeatSpec `json:"transport,omitempty"`
}

// AppHeartbeatSpec defines an application-level ping/pong protocol.
//
// Common patterns:
//
//	{"event":"ping"} -> reply {"event":"pong"}
//	{"op":"ping"}    -> reply {"op":"pong"}
//
// The host will auto-reply to inbound pings and can optionally send pings.
type AppHeartbeatSpec struct {
	MatchJSONField string `json:"matchJsonField"` // e.g. "event" or "op"
	PingValue      string `json:"pingValue"`      // e.g. "ping"
	PongValue      string `json:"pongValue"`      // e.g. "pong"

	// Optional periodic client ping to keep the connection alive.
	ClientPingIntervalMs int `json:"clientPingIntervalMs,omitempty"`
}

func (h AppHeartbeatSpec) Validate() error {
	if strings.TrimSpace(h.MatchJSONField) == "" {
		return fmt.Errorf("heartbeat.app.matchJsonField is required")
	}
	if strings.TrimSpace(h.PingValue) == "" {
		return fmt.Errorf("heartbeat.app.pingValue is required")
	}
	if strings.TrimSpace(h.PongValue) == "" {
		return fmt.Errorf("heartbeat.app.pongValue is required")
	}
	if h.ClientPingIntervalMs < 0 {
		return fmt.Errorf("heartbeat.app.clientPingIntervalMs must be >= 0")
	}
	return nil
}

// TransportHeartbeatSpec defines transport-level keepalive via WS ping frames.
// If set, the host will write WS Ping control frames periodically.
type TransportHeartbeatSpec struct {
	PingIntervalMs int `json:"pingIntervalMs,omitempty"`
}

func (h TransportHeartbeatSpec) Validate() error {
	if h.PingIntervalMs < 0 {
		return fmt.Errorf("heartbeat.transport.pingIntervalMs must be >= 0")
	}
	return nil
}
