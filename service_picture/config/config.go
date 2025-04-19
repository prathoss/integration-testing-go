package config

import (
	"log/slog"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel           slog.Level    `envconfig:"LOG_LEVEL"            default:"info"`
	ServerPort         uint16        `envconfig:"SERVER_PORT"          default:"8080"`
	ServerAddress      string        `envconfig:"SERVER_ADDRESS"       default:""`
	ServerReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT"  default:"20s"`
	ServerWriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"1s"`
	DatabaseURI        string        `envconfig:"DATABASE_URI"                        required:"true"`
}

func NewFromEnv() (Config, error) {
	var c Config
	err := envconfig.Process("gopic", &c)
	return c, err
}
