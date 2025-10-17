package plugin

import (
	"github.com/extism/go-pdk"
)

// ConfigField defines a configuration field that a plugin requires
// This is used to generate UI forms for setting up connections
type ConfigField struct {
	Name        string         `json:"name"`                  // Field name (e.g., "apiKey", "applicationID")
	Label       string         `json:"label"`                 // Human-readable label for UI
	Type        string         `json:"type"`                  // Input type: "text", "password", "number", etc.
	Required    bool           `json:"required"`              // Whether this field is mandatory
	Encrypt     bool           `json:"encrypt"`               // Whether to encrypt this field in database
	Mask        bool           `json:"mask"`                  // Whether to mask this field in API responses
	Placeholder string         `json:"placeholder,omitempty"` // Placeholder text for UI
	Description string         `json:"description,omitempty"` // Help text explaining the field
	Default     any            `json:"default,omitempty"`     // Default value
	Options     map[string]any `json:"options,omitempty"`     // Type-specific options
}

// ExportConfigFields exports configuration fields as JSON
func ExportConfigFields(fields []ConfigField) int32 {
	pdk.OutputJSON(fields)
	return 0
}

// ReadConfig reads configuration from plugin input (used in init export)
func ReadConfig() (map[string]any, error) {
	var config map[string]any
	err := pdk.InputJSON(&config)
	return config, err
}
