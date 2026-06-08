// Package config loads all the config values required for the app
package config

import (
	"errors"
	"os"
)

type Config struct {
	DBPath        string
	MigrationPath string
}

func Load() (Config, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	migrationPath := os.Getenv("DB_MIGRATION_PATH")
	config := Config{
		DBPath:        dbPath,
		MigrationPath: migrationPath,
	}

	var err error = nil
	if dbPath == "" || migrationPath == "" {
		err = errors.New("DATABASE_PATH and DB_MIGRATION_PATH is required")
	}

	return config, err
}
