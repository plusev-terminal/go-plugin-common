package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	mapstructure "github.com/go-viper/mapstructure/v2"
)

// MapToStruct populates a struct from a map and validates it.
// T is the type of the struct to populate (must be a pointer to a struct).
// Returns an error if parsing or validation fails.
func MapToStruct[T any](data map[string]any, target *T) error {
	// Initialize mapstructure decoder
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc(time.RFC3339Nano), // Handle time parsing
		),
		Metadata:         nil,
		Result:           target,
		TagName:          "mapstructure", // Use mapstructure tags for field mapping
		WeaklyTypedInput: true,           // Allow flexible type conversions (e.g., string to int)
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create mapstructure decoder: %w", err)
	}

	// Parse map into struct
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to parse map into struct: %w", err)
	}

	// Initialize validator
	validate := validator.New()
	if err := validate.Struct(target); err != nil {
		// Customize error message for validation
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// StructToMap converts a struct to a map[string]any
func StructToMap(input any, output *map[string]any) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, output)
}
