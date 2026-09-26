package config

import (
	"os"
)

type Config struct {
	Addr  string
	DBDSN string
}

// load the config
func Load() Config {
	// load from env
	addr := os.Getenv("KDOCS_ADDR")
	dsn := os.Getenv("KDOCS_DB_DSN")

	if addr == "" {
		addr = ":8080"
	}

	if dsn == "" {
		dsn = "postgres://admin:Incorrect@localhost:5432/kdocs?sslmode=disable"
	}

	return Config{
		Addr:  addr,
		DBDSN: dsn,
	}
}
