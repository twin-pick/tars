package src

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	omdbKey, err := loadEnv("OMDB_API_KEY")
	if err != nil {
		return nil, err
	}
	scrapperPort, err := loadEnv("SCRAPPER_PORT")
	if err != nil {
		return nil, err
	}
	exposedPort, err := loadEnv("EXPOSED_PORT")
	if err != nil {
		return nil, err
	}

	return &Config{
		OMDBApiKey:   omdbKey,
		ScrapperPort: scrapperPort,
		ExposedPort:  exposedPort,
	}, nil
}

func loadEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable not set: %s", key)
	}
	return value, nil
}
