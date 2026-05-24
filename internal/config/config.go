// Package config loads all the config values required for the app
package config

import (
	"os"
	"errors"
)

type Config struct {
	DBPath string
}

func Load() (Config, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	var err error = nil
	config := Config{
		DBPath: dbPath,
	}

	if dbPath == "" {
		err = errors.New("DATABASE_PATH is required")
	}

	return config, err
}

