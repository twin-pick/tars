package src

import (
	"os"
	"testing"
)

func TestNewConfig(T *testing.T) {
	os.Setenv("OMDB_API_KEY", "testkey")
	os.Setenv("SCRAPPER_PORT", "8080")
	os.Setenv("EXPOSED_PORT", "9090")

	config, err := NewConfig()
	if err != nil {
		T.Fatalf("Expected no error, got %v", err)
	}

	if config.OMDbApiKey != "testkey" {
		T.Errorf("Expected OMDBApiKey to be 'testkey', got %s", config.OMDbApiKey)
	}
	if config.ScrapperPort != "8080" {
		T.Errorf("Expected ScrapperPort to be '8080', got %s", config.ScrapperPort)
	}
	if config.ExposedPort != "9090" {
		T.Errorf("Expected ExposedPort to be '9090', got %s", config.ExposedPort)
	}
}
