package src

import (
	"os"
	"testing"
)

func TestNewConfig_Success(t *testing.T) {
	os.Setenv("OMDB_API_KEY", "dummy_key")
	os.Setenv("SCRAPPER_PORT", "3000")
	os.Setenv("EXPOSED_PORT", "8080")
	defer os.Clearenv()

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.OMDBApiKey != "dummy_key" {
		t.Errorf("Expected OMDBApiKey to be 'dummy_key', got %s", cfg.OMDBApiKey)
	}
	if cfg.ScrapperPort != "3000" {
		t.Errorf("Expected ScrapperPort to be '3000', got %s", cfg.ScrapperPort)
	}
	if cfg.ExposedPort != "8080" {
		t.Errorf("Expected ExposedPort to be '8080', got %s", cfg.ExposedPort)
	}
}

func TestNewConfig_MissingVar(t *testing.T) {
	os.Clearenv()

	_, err := NewConfig()
	if err == nil {
		t.Fatal("Expected error when env vars are missing, got nil")
	}
}

func TestLoadEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	expected := "value123"
	os.Setenv(key, expected)
	defer os.Unsetenv(key)

	val, err := loadEnv(key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val != expected {
		t.Errorf("Expected %s, got %s", expected, val)
	}
}

func TestLoadEnv_Missing(t *testing.T) {
	os.Unsetenv("MISSING_VAR")
	_, err := loadEnv("MISSING_VAR")
	if err == nil {
		t.Error("Expected error for missing env var, got nil")
	}
}
