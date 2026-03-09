package utils

import (
	"os"
	"strconv"
)

func GetEnvVarOrDefault[T any](key string, defaultValue T) T {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	var result T
	var err error

	switch any(defaultValue).(type) {
	case string:
		result = any(valueStr).(T)
	case int:
		var val int
		val, err = strconv.Atoi(valueStr)
		result = any(val).(T)
	case int64:
		var val int64
		val, err = strconv.ParseInt(valueStr, 10, 64)
		result = any(val).(T)
	case float64:
		var val float64
		val, err = strconv.ParseFloat(valueStr, 64)
		result = any(val).(T)
	case bool:
		var val bool
		val, err = strconv.ParseBool(valueStr)
		result = any(val).(T)
	default:
		// For unsupported types, return the default value
		return defaultValue
	}

	if err != nil {
		return defaultValue
	}

	return result
}
