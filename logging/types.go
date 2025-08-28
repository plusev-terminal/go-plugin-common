package logging

import (
	"time"
)

// PluginLogRecord represents a log record that will be sent to the host
type PluginLogRecord struct {
	ID        uint64         `json:"id"`
	PluginID  string         `json:"pluginId"`
	EventType string         `json:"eventType"`
	Timestamp time.Time      `json:"timestamp"`
	Message   string         `json:"message"`
	Data      map[string]any `json:"data,omitempty"`
}
