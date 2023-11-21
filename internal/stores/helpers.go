package stores

import (
	"encoding/json"
	"fmt"
)

// Helper function to convert value to JSON string
func toJSON(value interface{}) string {
	bytes, _ := json.Marshal(value)
	return string(bytes)
}

// Helper function to convert array of values to array of JSON strings
func toJSONArray(values []interface{}) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = toJSON(v)
	}
	return result
}

// Helper function to convert interface{} to string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		if v == nil {
			return "null"
		}
		return fmt.Sprintf("%v", val)
	}
}
