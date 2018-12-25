// +build release

package config

// Config ...
type Config struct {
	Debug     bool   `default:"false"`
	Port      string `default:"8080"`
	GrpcPort  string `default:"50151"`
	LogRotate LogRotate
}
