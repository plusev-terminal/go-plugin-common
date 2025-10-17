package plugin

import (
	"encoding/json"

	"github.com/extism/go-pdk"
)

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
	var config map[string]any
	err := pdk.InputJSON(&config)
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
