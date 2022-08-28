package config

import (
	"github.com/jkmrto/trade_executor/infra/http"
	"github.com/jkmrto/trade_executor/infra/sqlite3"
)

// Config is a container for all the needed app configuration.
type Config struct {
	HTTP    http.Config
	Sqlite3 sqlite3.Config
}

// New is a constructor.
func New() Config {
	return Config{
		HTTP: http.Config{
			Address: ":8080",
		},
		Sqlite3: sqlite3.Config{
			DatabaseFilePath:    "./dev.db",
			MigrationsDirectory: "./migrations",
		},
	}
}
