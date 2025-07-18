package src

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func NewConfig() *Config {
	godotenv.Load()

	return &Config{
		OMDBApiKey:   loadEnv("OMDB_API_KEY"),
		ScrapperPort: loadEnv("SCRAPPER_PORT"),
		ExposedPort:  loadEnv("EXPOSED_PORT"),
	}
}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable not set: %s", key)
	}
	return value
}
