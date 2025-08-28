package logging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/extism/go-pdk"
)

// Import the log_record host function
//
//go:wasmimport extism:host/user log_record
func hostLogRecord(offset uint64) uint64

// Logger provides logging functionality for plugins
type Logger struct {
	pluginID string
}

// NewLogger creates a new logger instance
// The pluginID will be automatically set by the host, but can be overridden if needed
func NewLogger(pluginID string) *Logger {
	return &Logger{
		pluginID: pluginID,
	}
}

// NewLogRecord creates a new log record with the current timestamp
func (l *Logger) NewLogRecord(eventType string) *PluginLogRecord {
	return &PluginLogRecord{
		PluginID:  l.pluginID,
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      make(map[string]any),
	}
}

// SetData sets the data field for the log record
func (r *PluginLogRecord) SetData(data map[string]any) *PluginLogRecord {
	r.Data = data
	return r
}

// SetMessage sets the message field for the log record
func (r *PluginLogRecord) SetMessage(message string) *PluginLogRecord {
	r.Message = message
	return r
}

// AddData adds a key-value pair to the data field
func (r *PluginLogRecord) AddData(key string, value any) *PluginLogRecord {
	if r.Data == nil {
		r.Data = make(map[string]any)
	}
	r.Data[key] = value
	return r
}

// Record sends the log record to the host via the log_record host function
func (r *PluginLogRecord) Record() error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
	}

	mem := pdk.AllocateBytes(data)
	defer mem.Free()

	// Call the log_record host function
	hostLogRecord(mem.Offset())

	return nil
}

// Convenience methods for common log levels

// Info logs an info message
func (l *Logger) Info(message string) error {
	return l.NewLogRecord("info").SetMessage(message).Record()
}

// InfoWithData logs an info message with additional data
func (l *Logger) InfoWithData(message string, data map[string]any) error {
	return l.NewLogRecord("info").SetMessage(message).SetData(data).Record()
}

// Error logs an error message
func (l *Logger) Error(message string) error {
	return l.NewLogRecord("error").SetMessage(message).Record()
}

// ErrorWithData logs an error message with additional data
func (l *Logger) ErrorWithData(message string, data map[string]any) error {
	return l.NewLogRecord("error").SetMessage(message).SetData(data).Record()
}

// Warn logs a warning message
func (l *Logger) Warn(message string) error {
	return l.NewLogRecord("warn").SetMessage(message).Record()
}

// WarnWithData logs a warning message with additional data
func (l *Logger) WarnWithData(message string, data map[string]any) error {
	return l.NewLogRecord("warn").SetMessage(message).SetData(data).Record()
}

// Debug logs a debug message
func (l *Logger) Debug(message string) error {
	return l.NewLogRecord("debug").SetMessage(message).Record()
}

// DebugWithData logs a debug message with additional data
func (l *Logger) DebugWithData(message string, data map[string]any) error {
	return l.NewLogRecord("debug").SetMessage(message).SetData(data).Record()
}

// Event logs a custom event
func (l *Logger) Event(eventType, message string) error {
	return l.NewLogRecord(eventType).SetMessage(message).Record()
}

// EventWithData logs a custom event with additional data
func (l *Logger) EventWithData(eventType, message string, data map[string]any) error {
	return l.NewLogRecord(eventType).SetMessage(message).SetData(data).Record()
}
