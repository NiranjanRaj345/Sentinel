package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	recoveryrepo "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/recovery"
)

type Store struct {
	db *sql.DB
}

func OpenRecovery(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("recovery storage path cannot be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open recovery database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS recovery_executions (
		id         TEXT    NOT NULL,
		policy_id  TEXT    NOT NULL,
		status     TEXT    NOT NULL,
		current    INTEGER NOT NULL,
		attempts   INTEGER NOT NULL,
		message    TEXT    NOT NULL,
		started_at DATETIME NOT NULL,
		finished_at DATETIME,
		PRIMARY KEY (id)
	)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create recovery executions table: %w", err)
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_recovery_executions_started_at ON recovery_executions(started_at)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create recovery executions index: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Save(ctx context.Context, execution recoveryrepo.Execution) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("recovery store is not open")
	}

	var finishedAt interface{}
	if execution.FinishedAt != nil {
		finishedAt = *execution.FinishedAt
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO recovery_executions (id, policy_id, status, current, attempts, message, started_at, finished_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET status=excluded.status, current=excluded.current, attempts=excluded.attempts, message=excluded.message, finished_at=excluded.finished_at`,
		execution.ID, execution.PolicyID, execution.Status, execution.Current, execution.Attempts, execution.Message, execution.StartedAt, finishedAt,
	)
	if err != nil {
		return fmt.Errorf("save recovery execution: %w", err)
	}

	return nil
}

func (s *Store) Recent(limit int) ([]recoveryrepo.Execution, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("recovery store is not open")
	}

	if limit <= 0 {
		limit = 100
	}

	rows, err := s.db.Query(
		`SELECT id, policy_id, status, current, attempts, message, started_at, finished_at FROM recovery_executions ORDER BY started_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query recent recovery executions: %w", err)
	}
	defer rows.Close()

	var executions []recoveryrepo.Execution

	for rows.Next() {
		var e recoveryrepo.Execution
		var finishedAt sql.NullTime

		if err := rows.Scan(&e.ID, &e.PolicyID, &e.Status, &e.Current, &e.Attempts, &e.Message, &e.StartedAt, &finishedAt); err != nil {
			return nil, fmt.Errorf("scan recovery execution: %w", err)
		}

		if finishedAt.Valid {
			finished := finishedAt.Time
			e.FinishedAt = &finished
		}

		executions = append(executions, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recovery executions: %w", err)
	}

	return executions, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
