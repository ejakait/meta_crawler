package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

func createTables(ctx context.Context, logger *slog.Logger, db *sql.DB) error {
	logger.InfoContext(ctx, "creating tables")
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS datasets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			owner TEXT,
			tags TEXT
		);
	`)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create table", "error", err)
		return fmt.Errorf("failed to create table: %w", err)
	}
	logger.InfoContext(ctx, "created datasets table")
	// _, err = db.ExecContext(ctx, `
	// 	CREATE TABLE IF NOT EXISTS metadata (
	// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 		link_id INTEGER NOT NULL REFERENCES links(id),
	// 		key TEXT NOT NULL,
	// 		value TEXT NOT NULL,
	// 		UNIQUE(link_id, key)
	// 	);
	// `)
	// if err != nil {
	// 	logger.ErrorContext(ctx, "failed to create metadata table", "error", err)
	// 	return fmt.Errorf("failed to create metadata table: %w", err)
	// }

	return nil
}

func InitDB(ctx context.Context, logger *slog.Logger, dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.ErrorContext(ctx, "failed to open database", "error", err)
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to close database", "error", err)
		}
	}()

	if err := createTables(ctx, logger, db); err != nil {
		logger.ErrorContext(ctx, "failed to instantiate tables", "error", err)
	}

	stmt := `select * from datasets`
	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		logger.ErrorContext(ctx, "failed to query datasets", "error", err)
		return fmt.Errorf("failed to query datasets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var description string
		var owner string
		var tags string
		if err := rows.Scan(&id, &name, &description, &owner, &tags); err != nil {
			logger.ErrorContext(ctx, "failed to scan dataset", "error", err)
			return fmt.Errorf("failed to scan dataset: %w", err)
		}
		logger.InfoContext(ctx, "dataset found", "id", id, "name", name, "description", description, "owner", owner, "tags", tags)
	}

	return nil
}
