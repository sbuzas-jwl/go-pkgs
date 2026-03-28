package main

import "github.com/sbuzas-jwl/go-pkgs/internal/sqlite"

const (
	// DefaultConfigPath is the default path to the application configuration.
	DefaultConfigPath = "/opt/todo/todod.conf"
)

type Config struct {
	HTTP struct {
		Port   int    `env:"HTTP_PORT"`
		Domain string `env:"HTTP_DOMAIN"`
	}
	DB sqlite.Config

	Debug string `env:"DEBUG"`
}

// DefaultConfig returns a new instance of Config with defaults set.
func DefaultConfig() Config {
	var config Config
	config.HTTP.Port = 8080
	return config
}
