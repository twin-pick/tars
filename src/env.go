package src

import (
	"fmt"
	"os"
)

func loadEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable not set: %s", key)
	}
	return value, nil
}
