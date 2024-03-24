package env

import (
	"fmt"
	"os"
	"strconv"
)

func String(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("env empty key: %v", key))
	}

	return val
}

func Int(key string) int {
	val := String(key)

	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("env int error: %v", err))
	}

	return intVal
}

func Float64(key string) float64 {
	val := String(key)

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(fmt.Sprintf("env float error: %v", err))
	}

	return floatVal
}

func Bool(key string) bool {
	val := String(key)

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		panic(fmt.Sprintf("env bool error: %v", err))
	}

	return boolVal
}
