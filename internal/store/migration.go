package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func RunMigration(database *sql.DB, migrationsDir string) error {
	if _, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`); err != nil {
		return fmt.Errorf("create schema_migrations tabe: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("find migration files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationsDir)
	}

	sort.Strings(files)

	for _, file := range files {
		version := filepath.Base(file)

		applied, err := migrationApplied(database, version)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", version, err)
		}

		tx, err := database.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", version, err)
		}
		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("execute migration %s: %w", version, err)
		}
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", version, err)
		}
	}
	return nil
}

func migrationApplied(database *sql.DB, version string) (bool, error) {
	var count int

	err := database.QueryRow(
		"SELECT COUNT(*) FROM schema_migrations WHERE version = ?",
		version,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}

	return count > 0, nil
}
