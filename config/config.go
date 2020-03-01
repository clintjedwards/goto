package config

import "github.com/kelseyhightower/envconfig"

// Config refers to general application configuration
type Config struct {
	Debug       bool   `envconfig:"debug" default:"false"`
	LogLevel    string `envconfig:"loglevel" default:"info"`
	Host        string `envconfig:"host" default:"localhost:8080"`
	MaxIDLength int    `envconfig:"max_id_length" default:"50"` // The total amount of characters that a short name can be
	Database    *DatabaseConfig
}

// BoltConfig represents a on-disk key/value store
// https://github.com/boltdb/bolt
type BoltConfig struct {
	// file path for database file
	Path string `envconfig:"database_path_bolt" default:"/tmp/go.db"`
}

// RedisConfig represents a key/value store
// https://redis.io
type RedisConfig struct {
	Host     string `envconfig:"database_host_redis" default:"localhost:6379"`
	Password string `envconfig:"database_password_redis"`
	DB       int    `envconfig:"database_db_redis" default:"0"` // redis database number 0-15
}

// DatabaseConfig defines config settings for comet database
type DatabaseConfig struct {
	// The database engine used by the backend
	// possible values are: bolt, redis
	Engine string `envconfig:"database_engine" default:"bolt"`
	Bolt   *BoltConfig
	Redis  *RedisConfig
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
