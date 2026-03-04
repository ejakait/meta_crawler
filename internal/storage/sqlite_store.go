package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

func createTables(ctx context.Context, db *sql.DB) error {
	slog.InfoContext(ctx, "creating tables")
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS datasets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT NULL,
			owner TEXT NULL,
			object_count INTEGER NULL,
			location TEXT NULL,
			tags TEXT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create table", "error", err)
		return fmt.Errorf("failed to create table: %w", err)
	}
	slog.InfoContext(ctx, "created datasets table")
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
	// 	slog.ErrorContext(ctx, "failed to create metadata table", "error", err)
	// 	return fmt.Errorf("failed to create metadata table: %w", err)
	// }

	return nil
}

func InitDB(ctx context.Context, dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.ErrorContext(ctx, "failed to open database", "error", err)
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to close database", "error", err)
		}
	}()

	if err := createTables(ctx, db); err != nil {
		slog.ErrorContext(ctx, "failed to instantiate tables", "error", err)
	}

	// stmt := `select * from datasets`
	// rows, err := db.QueryContext(ctx, stmt)
	// if err != nil {
	// 	slog.ErrorContext(ctx, "failed to query datasets", "error", err)
	// 	return fmt.Errorf("failed to query datasets: %w", err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	var description string
	// 	var owner string
	// 	var tags string
	// 	if err := rows.Scan(&id, &name, &description, &owner, &tags); err != nil {
	// 		slog.ErrorContext(ctx, "failed to scan dataset", "error", err)
	// 		return fmt.Errorf("failed to scan dataset: %w", err)
	// 	}
	// 	slog.InfoContext(ctx, "dataset found", "id", id, "name", name, "description", description, "owner", owner, "tags", tags)
	// }

	return nil
}

func OpenDB(ctx context.Context, dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.ErrorContext(ctx, "failed to open database", "error", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func SaveContainerMetadata(ctx context.Context, db *sql.DB, containerMetadata ContainerMetadata) error {
	stmt := `INSERT INTO datasets (name, description, owner, object_count, location, tags) VALUES (?, ?, ?, ?, ?, ?) on conflict (name) do update set description = ?, owner = ?, object_count = ?, location = ?, tags = ?, updated_at = CURRENT_TIMESTAMP`
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			slog.ErrorContext(ctx, "failed to rollback transaction", "error", err)
		}
	}()

	_, err = tx.ExecContext(ctx, stmt, containerMetadata.Name, containerMetadata.Description, containerMetadata.Owner, containerMetadata.ObjectCount, containerMetadata.Location, containerMetadata.Tags, containerMetadata.Description, containerMetadata.Owner, containerMetadata.ObjectCount, containerMetadata.Location, containerMetadata.Tags)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save container metadata", "error", err)
		return fmt.Errorf("failed to save container metadata: %w", err)
	}

	if err := tx.Commit(); err != nil {
		slog.ErrorContext(ctx, "failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
