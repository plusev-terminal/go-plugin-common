package datapipe

import (
	"encoding/json"

	"github.com/extism/go-pdk"
	dt "github.com/plusev-terminal/go-plugin-common/datapipe/types"
)

// ReadProcessRequest reads a process request from plugin input
func ReadProcessRequest() (dt.ProcessRequest, error) {
	var req dt.ProcessRequest
	err := pdk.InputJSON(&req)
	return req, err
}

// WriteProcessResponse writes a process response to plugin output
func WriteProcessResponse(resp dt.ProcessResponse) int32 {
	pdk.OutputJSON(resp)
	if resp.Success {
		return 0
	}
	return 1
}

// ReadConfig reads configuration from plugin input (used in init export)
func ReadConfig() (map[string]any, error) {
	var config map[string]any
	err := pdk.InputJSON(&config)
	return config, err
}

// SuccessResponse creates a successful process response
func SuccessResponse(output map[string]any) dt.ProcessResponse {
	return dt.ProcessResponse{
		Success: true,
		Output:  output,
	}
}

// ErrorResponse creates an error process response
func ErrorResponse(err error) dt.ProcessResponse {
	return dt.ProcessResponse{
		Success: false,
		Error:   err.Error(),
	}
}

// ErrorResponseMsg creates an error process response with a message
func ErrorResponseMsg(msg string) dt.ProcessResponse {
	return dt.ProcessResponse{
		Success: false,
		Error:   msg,
	}
}

// ConfigStore helps manage plugin configuration
type ConfigStore struct {
	config map[string]any
}

// NewConfigStore creates a new configuration store
func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		config: make(map[string]any),
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

// GetString retrieves a configuration value as string
func (cs *ConfigStore) GetString(key string) string {
	if val, ok := cs.config[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetStringOr retrieves a configuration value with a default
func (cs *ConfigStore) GetStringOr(key, defaultValue string) string {
	if val := cs.GetString(key); val != "" {
		return val
	}
	return defaultValue
}

// GetNumber retrieves a configuration value as float64
func (cs *ConfigStore) GetNumber(key string) float64 {
	if val, ok := cs.config[key]; ok {
		if num, ok := val.(float64); ok {
			return num
		}
	}
	return 0
}

// GetNumberOr retrieves a configuration value with a default
func (cs *ConfigStore) GetNumberOr(key string, defaultValue float64) float64 {
	if val := cs.GetNumber(key); val != 0 {
		return val
	}
	return defaultValue
}

// GetBool retrieves a configuration value as bool
func (cs *ConfigStore) GetBool(key string) bool {
	if val, ok := cs.config[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// Get retrieves a configuration value (raw interface{})
func (cs *ConfigStore) Get(key string) any {
	return cs.config[key]
}

// Has checks if a configuration key exists
func (cs *ConfigStore) Has(key string) bool {
	_, ok := cs.config[key]
	return ok
}
