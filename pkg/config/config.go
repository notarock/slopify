package config

import (
	"os"
)

// Config represents the configuration parsed from the .env file
type Config struct {
	OpenaiKey string
}

func GetEnvConfig() (Config, error) {
	config := Config{}

	// Parse environment variables
	config.OpenaiKey = os.Getenv("OPENAI_KEY")

	return config, nil
}
