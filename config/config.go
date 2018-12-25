package config

import "github.com/kelseyhightower/envconfig"

// New ...
func New() Config {
	config := Config{}
	envconfig.MustProcess("", &config)

	return config
}

// LogRotate setting
type LogRotate struct {
	Filename   string `default:"logs/main.log"`
	MaxSize    int    `default:"10"` // megabytes
	MaxBackups int    `default:"10"`
	MaxAge     int    `default:"365"`  //days
	Compress   bool   `default:"true"` // disabled by default
}
