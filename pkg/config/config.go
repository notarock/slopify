package config

import (
	"os"
)

// Config represents the configuration parsed from the .env file
type Config struct {
	PexelsAPIKey string
}

func GetEnvConfig() (Config, error) {
	config := Config{}

	// Parse environment variables
	config.PexelsAPIKey = os.Getenv("PEXELS_API_KEY")

	return config, nil
}
