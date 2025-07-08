package tars

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func NewConfig() Config {
	godotenv.Load()

	cfg := Config{
		OMDBApiKey:   os.Getenv("OMDB_API_KEY"),
		ScrapperPort: os.Getenv("SCRAPPER_PORT"),
		ExposedPort:  os.Getenv("EXPOSED_PORT"),
	}

	if cfg.OMDBApiKey == "" {
		log.Fatal("TMDB_TOKEN env var is not set")
	}
	if cfg.ScrapperPort == "" {
		log.Fatal("SCRAPPER_PORT env var is not set")
	}
	if cfg.ExposedPort == "" {
		log.Fatal("EXPOSED_PORT env var is not set")
	}

	return cfg
}
