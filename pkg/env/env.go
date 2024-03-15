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
