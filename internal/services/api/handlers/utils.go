package handlers

import (
	"strconv"
)

// parseID safely parses an ID from various types (string, float64, int, etc.)
func parseID(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int64(val)
	case string:
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	case int:
		return int64(val)
	case int64:
		return val
	case int32:
		return int64(val)
	default:
		return 0
	}
}
