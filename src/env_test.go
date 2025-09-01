package src

import (
	"os"
	"testing"
)

func TestLoadEnv(T *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")
	value, err := loadEnv("TEST_ENV_VAR")
	if err != nil {
		T.Errorf("Expected no error, got %v", err)
	}
	if value != "test_value" {
		T.Errorf("Expected 'test_value', got %s", value)
	}

	_, err = loadEnv("NON_EXISTENT_VAR")
	if err == nil {
		T.Error("Expected an error for non-existent variable, got nil")
	}
}
