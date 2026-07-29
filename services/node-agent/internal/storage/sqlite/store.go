package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("storage path cannot be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := createSchema(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *Store) Save(info metrics.Info) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is not open")
	}

	payload, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal metrics: %w", err)
	}

	_, err = s.db.Exec(
		`INSERT INTO metrics_snapshots (timestamp, payload) VALUES (?, ?)`,
		info.Metadata.Timestamp,
		payload,
	)
	if err != nil {
		return fmt.Errorf("insert snapshot: %w", err)
	}

	return nil
}
