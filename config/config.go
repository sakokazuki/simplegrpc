package config

import "github.com/kelseyhightower/envconfig"

// New ...
func New() Config {
	config := Config{}
	envconfig.MustProcess("", &config)

	return config
}

// Config ...
type Config struct {
	Debug    bool   `default:"true"`
	Port     string `default:"8080"`
	GrpcPort string `default:"50151"`
}
