package datasrc

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
	dt "github.com/plusev-terminal/go-plugin-common/datasrc/types"
)

// ================== Helper Functions for Plugin Exports ==================

// ExportConfigFields exports configuration fields as JSON
func ExportConfigFields(fields []dt.ConfigField) int32 {
	pdk.OutputJSON(fields)
	return 0
}

// ReadConfig reads configuration from plugin input (used in init export)
func ReadConfig() (map[string]string, error) {
	var config map[string]string
	err := pdk.InputJSON(&config)
	return config, err
}

// ReadCommand reads a command from plugin input (used in handle_command export)
func ReadCommand() (dt.Command, error) {
	var cmd dt.Command
	err := pdk.InputJSON(&cmd)
	return cmd, err
}

// WriteResponse writes a response to plugin output
func WriteResponse(resp dt.Response) int32 {
	pdk.OutputJSON(resp)
	if resp.Success {
		return 0
	}
	return 1
}

// SuccessResponse creates a successful response with data
func SuccessResponse(data any) dt.Response {
	return dt.Response{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(err error) dt.Response {
	return dt.Response{
		Success: false,
		Error:   err.Error(),
	}
}

// ErrorResponseMsg creates an error response with a message
func ErrorResponseMsg(msg string) dt.Response {
	return dt.Response{
		Success: false,
		Error:   msg,
	}
}

// Common timeframes that most exchanges support
var CommonTimeframes = []dt.Timeframe{
	{Label: "1 minute", Value: "1m", Interval: 60},
	{Label: "5 minutes", Value: "5m", Interval: 300},
	{Label: "15 minutes", Value: "15m", Interval: 900},
	{Label: "30 minutes", Value: "30m", Interval: 1800},
	{Label: "1 hour", Value: "1h", Interval: 3600},
	{Label: "4 hours", Value: "4h", Interval: 14400},
	{Label: "1 day", Value: "1d", Interval: 86400},
}

// ================== Command Router Helper ==================

// CommandHandler is a function that handles a specific command
type CommandHandler func(params map[string]any) dt.Response

// CommandRouter helps route commands to handlers
type CommandRouter struct {
	handlers map[string]CommandHandler
}

// NewCommandRouter creates a new command router
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		handlers: make(map[string]CommandHandler),
	}
}

// Register registers a handler for a command name
func (r *CommandRouter) Register(commandName string, handler CommandHandler) {
	r.handlers[commandName] = handler
}

// Handle routes a command to the appropriate handler
func (r *CommandRouter) Handle(cmd dt.Command) dt.Response {
	handler, ok := r.handlers[cmd.Name]
	if !ok {
		return ErrorResponseMsg(fmt.Sprintf("unknown command: %s", cmd.Name))
	}
	return handler(cmd.Params)
}

// HandleJSON reads command from input, routes it, and writes response
func (r *CommandRouter) HandleJSON() int32 {
	cmd, err := ReadCommand()
	if err != nil {
		return WriteResponse(ErrorResponse(err))
	}
	return WriteResponse(r.Handle(cmd))
}

// ================== Configuration Storage Helper ==================

// ConfigStore helps manage plugin configuration
type ConfigStore struct {
	config map[string]string
}

// NewConfigStore creates a new configuration store
func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		config: make(map[string]string),
	}
}

// Load loads configuration from JSON input
func (cs *ConfigStore) Load() error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}
	cs.config = config
	return nil
}

// LoadFromBytes loads configuration from JSON bytes
func (cs *ConfigStore) LoadFromBytes(data []byte) error {
	return json.Unmarshal(data, &cs.config)
}

// Get retrieves a configuration value
func (cs *ConfigStore) Get(key string) string {
	return cs.config[key]
}

// GetOr retrieves a configuration value with a default
func (cs *ConfigStore) GetOr(key, defaultValue string) string {
	if val, ok := cs.config[key]; ok && val != "" {
		return val
	}
	return defaultValue
}

// Has checks if a configuration key exists
func (cs *ConfigStore) Has(key string) bool {
	_, ok := cs.config[key]
	return ok
}
