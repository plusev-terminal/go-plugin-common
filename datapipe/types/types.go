package types

// DataType represents the type of data flowing through node ports
type DataType string

const (
	DataTypeOHLCVRecord DataType = "OHLCVRecord"
	DataTypeSignal      DataType = "Signal"
	DataTypeStartSignal DataType = "StartSignal"
)

// NodePort defines an input or output port on a node
type NodePort struct {
	Name      string     `json:"name"`
	DataTypes []DataType `json:"dataTypes"`
}

// ConfigField defines a configuration field for the plugin
type ConfigField struct {
	Label       string         `json:"label"`       // Human-readable label
	Name        string         `json:"name"`        // Field key in config map
	Type        string         `json:"type"`        // Field type (string, number, boolean, select, etc.)
	Required    bool           `json:"required"`    // Whether the field is required
	Default     any            `json:"default"`     // Default value
	Description string         `json:"description"` // Help text
	Options     map[string]any `json:"options"`     // Type-specific options
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

// GuiDefinition defines the configuration UI for the plugin
type GuiDefinition struct {
	Controls []GuiControl `json:"controls"`
}

// GuiControl defines a single UI control in the configuration GUI
type GuiControl struct {
	Label   string         `json:"label"`
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Options map[string]any `json:"options"`
}

// OHLCVRecord represents a single candlestick/OHLCV data point
type OHLCVRecord struct {
	Timestamp int64   `json:"timestamp"` // Unix timestamp in milliseconds
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

// ToMap converts OHLCVRecord to a map for processing
func (o *OHLCVRecord) ToMap() map[string]any {
	return map[string]any{
		"timestamp": o.Timestamp,
		"open":      o.Open,
		"high":      o.High,
		"low":       o.Low,
		"close":     o.Close,
		"volume":    o.Volume,
	}
}

// FromMap populates OHLCVRecord from a map
func (o *OHLCVRecord) FromMap(m map[string]any) {
	if v, ok := m["timestamp"].(float64); ok {
		o.Timestamp = int64(v)
	}
	if v, ok := m["open"].(float64); ok {
		o.Open = v
	}
	if v, ok := m["high"].(float64); ok {
		o.High = v
	}
	if v, ok := m["low"].(float64); ok {
		o.Low = v
	}
	if v, ok := m["close"].(float64); ok {
		o.Close = v
	}
	if v, ok := m["volume"].(float64); ok {
		o.Volume = v
	}
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
