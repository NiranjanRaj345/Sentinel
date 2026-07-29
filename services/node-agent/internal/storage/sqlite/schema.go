package sqlite

import (
	"database/sql"
)

func createSchema(db *sql.DB) error {
	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		return err
	}

	metricsTable := `
CREATE TABLE IF NOT EXISTS metrics_snapshots (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp  DATETIME NOT NULL,
    payload    JSON     NOT NULL
)`

	if _, err := db.Exec(metricsTable); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL)`); err != nil {
		return err
	}

	var current int
	if err := db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&current); err != nil {
		return err
	}

	if current < SchemaVersion {
		if _, err := db.Exec(`INSERT INTO schema_version (version) VALUES (?)`, SchemaVersion); err != nil {
			return err
		}
	}

	return nil
}
