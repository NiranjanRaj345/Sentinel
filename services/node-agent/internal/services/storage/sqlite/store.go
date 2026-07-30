package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("service storage path cannot be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open services database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS services (
		name     TEXT    NOT NULL PRIMARY KEY,
		status   TEXT    NOT NULL,
		payload  JSON    NOT NULL
	)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create services table: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Save(ctx context.Context, item services.ServiceItem) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("services store is not open")
	}

	payload, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("marshal service item: %w", err)
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO services (name, status, payload) VALUES (?, ?, ?)`,
		item.Name, string(item.Status), payload,
	)
	if err != nil {
		return fmt.Errorf("insert service: %w", err)
	}
	return nil
}

func (s *Store) List(ctx context.Context) ([]services.ServiceItem, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("services store is not open")
	}

	rows, err := s.db.QueryContext(ctx, `SELECT name, status, payload FROM services`)
	if err != nil {
		return nil, fmt.Errorf("query services: %w", err)
	}
	defer rows.Close()

	var items []services.ServiceItem
	for rows.Next() {
		var name, status, payload string
		if err := rows.Scan(&name, &status, &payload); err != nil {
			return nil, fmt.Errorf("scan service: %w", err)
		}
		var item services.ServiceItem
		if err := json.Unmarshal([]byte(payload), &item); err != nil {
			return nil, fmt.Errorf("decode service item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate services: %w", err)
	}
	return items, nil
}

