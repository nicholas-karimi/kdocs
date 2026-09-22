package config

import (
	"os"
)

type Config struct {
	Addr string
}

// load the config
func Load() Config {
	// load from env
	addr := os.Getenv("KDOCS_ADDR")

	if addr == "" {
		addr = ":8080"
	}

	return Config{
		Addr: addr,
	}
}
