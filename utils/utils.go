package utils

import "os"

func GetEnvOrDefault(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultValue
}
