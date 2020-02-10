package config

import "github.com/kelseyhightower/envconfig"

// Config refers to general application configuration
type Config struct {
	Debug    bool   `envconfig:"debug" default:"false"`
	LogLevel string `envconfig:"loglevel" default:"info"`
	URL      string `envconfig:"url" default:"localhost:8080"`
}

// FromEnv pulls configration from environment variables
func FromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
