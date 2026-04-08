package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/armando284/tkt/internal/config"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	dbPath := config.GetDBPath()

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Better performance and foreign key support
	_, err = DB.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;")
	if err != nil {
		return err
	}

	return createTables()
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		root_path TEXT UNIQUE NOT NULL,
		name TEXT,
		registered_at TEXT DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		status TEXT DEFAULT 'todo',
		folder TEXT,
		branch TEXT,
		project_root TEXT NOT NULL,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(title, project_root),
		FOREIGN KEY (project_root) REFERENCES projects(root_path) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ticket_id INTEGER NOT NULL,
		start_ts TEXT NOT NULL,
		end_ts TEXT,
		duration INTEGER,
		FOREIGN KEY (ticket_id) REFERENCES tickets(id)
	);
	`

	_, err := DB.Exec(schema)
	return err
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
