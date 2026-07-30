package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"

	eventrepo "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
)

type Store struct {
	db *sql.DB
}

func OpenEvents(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("events storage path cannot be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open events database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS events (
		id         TEXT    NOT NULL,
		type       TEXT    NOT NULL,
		severity   TEXT    NOT NULL,
		source     TEXT    NOT NULL,
		title      TEXT    NOT NULL,
		message    TEXT    NOT NULL,
		metadata   JSON,
		created_at DATETIME NOT NULL
	)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create events table: %w", err)
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create events index: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Save(ctx context.Context, event eventrepo.Event) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("events store is not open")
	}

	payload, err := json.Marshal(event.Metadata)
	if err != nil {
		return fmt.Errorf("marshal event metadata: %w", err)
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO events (id, type, severity, source, title, message, metadata, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.Type, event.Severity, event.Source, event.Title, event.Message, payload, event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	return nil
}

func (s *Store) Recent(limit int) ([]eventrepo.Event, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("events store is not open")
	}

	if limit <= 0 {
		limit = 100
	}

	rows, err := s.db.Query(
		`SELECT id, type, severity, source, title, message, metadata, created_at FROM events ORDER BY created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query recent events: %w", err)
	}
	defer rows.Close()

	var eventsList []eventrepo.Event

	for rows.Next() {
		var e eventrepo.Event
		var payload []byte

		if err := rows.Scan(&e.ID, &e.Type, &e.Severity, &e.Source, &e.Title, &e.Message, &payload, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}

		if len(payload) > 0 {
			if err := json.Unmarshal(payload, &e.Metadata); err != nil {
				return nil, fmt.Errorf("decode event metadata: %w", err)
			}
		} else {
			e.Metadata = make(map[string]interface{})
		}

		eventsList = append(eventsList, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}

	return eventsList, nil
}
