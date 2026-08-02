package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("resource storage path cannot be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open resources database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS resources (
		id         TEXT    NOT NULL PRIMARY KEY,
		name       TEXT    NOT NULL,
		type       TEXT    NOT NULL,
		health     TEXT    NOT NULL,
		version    TEXT    NOT NULL,
		description TEXT   NOT NULL,
		provider   TEXT    NOT NULL,
		status     TEXT    NOT NULL,
		message    TEXT    NOT NULL
	)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create resources table: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Save(ctx context.Context, resource resources.Resource) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("resources store is not open")
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO resources (id, name, type, health, version, description, provider, status, message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		resource.ID, resource.Name, string(resource.Type), string(resource.Health), resource.Version, resource.Description, resource.Provider, resource.Status, resource.Message,
	)
	if err != nil {
		return fmt.Errorf("insert resource: %w", err)
	}
	return nil
}

func (s *Store) List(ctx context.Context) ([]resources.Resource, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("resources store is not open")
	}

	rows, err := s.db.QueryContext(ctx, `SELECT id, name, type, health, version, description, provider, status, message FROM resources`)
	if err != nil {
		return nil, fmt.Errorf("query resources: %w", err)
	}
	defer rows.Close()

	var items []resources.Resource
	for rows.Next() {
		var r resources.Resource
		if err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Health, &r.Version, &r.Description, &r.Provider, &r.Status, &r.Message); err != nil {
			return nil, fmt.Errorf("scan resource: %w", err)
		}
		items = append(items, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resources: %w", err)
	}
	return items, nil
}
