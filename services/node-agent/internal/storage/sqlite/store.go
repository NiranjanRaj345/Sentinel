package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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

func (s *Store) Latest() (metrics.Info, error) {
	if s == nil || s.db == nil {
		return metrics.Info{}, fmt.Errorf("store is not open")
	}

	row := s.db.QueryRow(
		`SELECT payload FROM metrics_snapshots ORDER BY timestamp DESC LIMIT 1`,
	)

	var payload []byte
	if err := row.Scan(&payload); err != nil {
		if err == sql.ErrNoRows {
			return metrics.Info{}, fmt.Errorf("no historical metrics available")
		}
		return metrics.Info{}, fmt.Errorf("query latest snapshot: %w", err)
	}

	var info metrics.Info
	if err := json.Unmarshal(payload, &info); err != nil {
		return metrics.Info{}, fmt.Errorf("decode snapshot: %w", err)
	}

	return info, nil
}

func (s *Store) Range(from, to time.Time) ([]metrics.Info, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("store is not open")
	}

	rows, err := s.db.Query(
		`SELECT payload FROM metrics_snapshots WHERE timestamp >= ? AND timestamp <= ? ORDER BY timestamp ASC`,
		from,
		to,
	)
	if err != nil {
		return nil, fmt.Errorf("query history range: %w", err)
	}
	defer rows.Close()

	var snapshots []metrics.Info

	for rows.Next() {
		var payload []byte
		if err := rows.Scan(&payload); err != nil {
			return nil, fmt.Errorf("scan snapshot: %w", err)
		}

		var info metrics.Info
		if err := json.Unmarshal(payload, &info); err != nil {
			return nil, fmt.Errorf("decode snapshot: %w", err)
		}

		snapshots = append(snapshots, info)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate snapshots: %w", err)
	}

	return snapshots, nil
}
