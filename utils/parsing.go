package utils

import "time"

type mapValue interface {
	float64 | string | bool
}

func IfThen[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func GetValue[T mapValue](key string, data map[string]any, defaultValue ...T) T {
	value, ok := data[key].(T)
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		var defaultValue T
		return defaultValue
	}

	var zero T

	if len(defaultValue) > 0 && value == zero {
		return defaultValue[0]
	}

	return value
}

func ExtractMap(key string, data map[string]any) map[string]any {
	if val, ok := data[key]; ok {
		if subMap, ok := val.(map[string]any); ok {
			return subMap
		}
	}

	return nil
}

func AnyMatches[T comparable](predicate func(T) bool, values ...T) bool {
	for _, v := range values {
		if predicate(v) {
			return true
		}
	}
	return false
}

// ExtractInt safely extracts an int value from the map
func ExtractInt(key string, data map[string]any) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return 0
}

// ExtractTime safely extracts a time.Time value from the map
// Supports: string (RFC3339), time.Time, int64/float64 (unix millis)
func ExtractTime(key string, data map[string]any) *time.Time {
	if val, ok := data[key]; ok && val != nil {
		switch v := val.(type) {
		case string:
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				return &t
			}
			if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
				return &t
			}
		case time.Time:
			return &v
		case int64:
			t := time.UnixMilli(v)
			return &t
		case float64:
			t := time.UnixMilli(int64(v))
			return &t
		}
	}
	return nil
}
