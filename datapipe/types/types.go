package types

// DataType represents the type of data flowing through node ports
type DataType string

const (
	DataTypeOHLCVRecord DataType = "OHLCVRecord"
	DataTypeSignal      DataType = "Signal"
	DataTypeStartSignal DataType = "StartSignal"
)

type NodeMeta struct {
	Name          string               `json:"name"`
	GuiDefinition GuiDefinition        `json:"guiDefinition"`
	Connections   Connections          `json:"connections"`
	CustomTypes   []DataTypeDefinition `json:"customTypes,omitempty"` // Plugin-defined data types with inheritance support
}

// DataTypeDefinition defines a custom data type with optional inheritance
type DataTypeDefinition struct {
	Name        string          `json:"name"`
	DisplayName string          `json:"displayName"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Extends     string          `json:"extends,omitempty"` // Inheritance: type this extends
	Fields      []DataTypeField `json:"fields"`
	Icon        string          `json:"icon,omitempty"`
	Color       string          `json:"color,omitempty"`
}

// DataTypeField represents a field within a data type
type DataTypeField struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"` // "string", "number", "timestamp", "boolean"
	Description string `json:"description"`
	Placeholder string `json:"placeholder"`
	Format      string `json:"format,omitempty"`
	Required    bool   `json:"required"`
}

type Connections struct {
	Inputs  []NodePort `json:"inputs"`
	Outputs []NodePort `json:"outputs"`
}

// NodePort defines an input or output port on a node
type NodePort struct {
	Name      string     `json:"name"`
	DataTypes []DataType `json:"dataTypes"`
}

// ProcessRequest contains the input data and configuration for processing
type ProcessRequest struct {
	Input  map[string]any `json:"input"`  // Input data keyed by port name
	Config map[string]any `json:"config"` // Plugin configuration values
}

// ProcessResponse contains the output data from processing
type ProcessResponse struct {
	Success bool           `json:"success"`          // Whether processing succeeded
	Output  map[string]any `json:"output,omitempty"` // Output data keyed by port name
	Error   string         `json:"error,omitempty"`  // Error message if failed
}

type GuiControlType string

const (
	TEXT_INPUT   GuiControlType = "text_input"
	NUMBER_INPUT GuiControlType = "number_input"
	CHECKBOX     GuiControlType = "checkbox"
	SELECT       GuiControlType = "select"
	MULTISELECT  GuiControlType = "multiselect"
)

// GuiDefinition defines the configuration UI for the plugin
type GuiDefinition struct {
	Controls []*GuiControl `json:"controls"`
}

// GuiControl defines a single UI control in the configuration GUI
type GuiControl struct {
	Label   string         `json:"label"`
	Name    string         `json:"name"`
	Type    GuiControlType `json:"type"`
	Options map[string]any `json:"options"`
}

// Signal represents a trading signal
type Signal struct {
	Type      string  `json:"type"`      // Signal type (buy, sell, etc.)
	Strength  float64 `json:"strength"`  // Signal strength (0-1)
	Timestamp int64   `json:"timestamp"` // Unix timestamp in milliseconds
	Message   string  `json:"message"`   // Optional message
}

// ToMap converts Signal to a map for processing
func (s *Signal) ToMap() map[string]any {
	return map[string]any{
		"type":      s.Type,
		"strength":  s.Strength,
		"timestamp": s.Timestamp,
		"message":   s.Message,
	}
}

// FromMap populates Signal from a map
func (s *Signal) FromMap(m map[string]any) {
	if v, ok := m["type"].(string); ok {
		s.Type = v
	}
	if v, ok := m["strength"].(float64); ok {
		s.Strength = v
	}
	if v, ok := m["timestamp"].(float64); ok {
		s.Timestamp = int64(v)
	}
	if v, ok := m["message"].(string); ok {
		s.Message = v
	}
}

// StartSignal represents a pipeline start signal
type StartSignal struct {
	Timestamp int64 `json:"timestamp"` // Unix timestamp in milliseconds
}

// ToMap converts StartSignal to a map for processing
func (ss *StartSignal) ToMap() map[string]any {
	return map[string]any{
		"timestamp": ss.Timestamp,
	}
}

// FromMap populates StartSignal from a map
func (ss *StartSignal) FromMap(m map[string]any) {
	if v, ok := m["timestamp"].(float64); ok {
		ss.Timestamp = int64(v)
	}
}
