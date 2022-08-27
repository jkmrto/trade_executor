package config

import (
	"github.com/jkmrto/trade_executor/infra/http"
)

// Config is a container for all the needed app configuration.
type Config struct {
	HTTP http.Config
}

// New is a constructor.
func New() Config {
	return Config{
		HTTP: http.Config{
			Address: ":8080",
		},
	}
}
