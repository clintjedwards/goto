package config

import "github.com/kelseyhightower/envconfig"

// Config refers to general application configuration
type Config struct {
	Debug         bool   `envconfig:"debug" default:"false"`
	LogLevel      string `envconfig:"loglevel" default:"info"`
	URL           string `envconfig:"url" default:"localhost:8080"`
	MaxNameLength int    `envconfig:"max_name_length" default:"50"` // The total amount of characters that a short name can be
	DBPath        string `envconfig:"db_path" default:"/tmp/go.db"`
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
