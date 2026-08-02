package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/automation"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("automation storage path cannot be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open automation database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS automation_executions (
		id         TEXT    NOT NULL PRIMARY KEY,
		rule_id    TEXT    NOT NULL,
		rule_name  TEXT    NOT NULL,
		action     TEXT    NOT NULL,
		success    INTEGER NOT NULL,
		message    TEXT    NOT NULL,
		created_at DATETIME NOT NULL
	)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create automation_executions table: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Save(ctx context.Context, record automation.ExecutionRecord) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("automation store is not open")
	}

	id := fmt.Sprintf("auto-%d", time.Now().UTC().UnixNano())
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO automation_executions (id, rule_id, rule_name, action, success, message, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, record.RuleID, record.RuleName, record.Action, boolToInt(record.Success), record.Message, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert automation execution: %w", err)
	}
	return nil
}

func (s *Store) List(ctx context.Context, limit int) ([]automation.ExecutionRecord, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("automation store is not open")
	}

	if limit <= 0 {
		limit = 100
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, rule_id, rule_name, action, success, message, created_at FROM automation_executions ORDER BY created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query automation executions: %w", err)
	}
	defer rows.Close()

	var records []automation.ExecutionRecord
	for rows.Next() {
		var r automation.ExecutionRecord
		var successInt int
		if err := rows.Scan(&r.ID, &r.RuleID, &r.RuleName, &r.Action, &successInt, &r.Message, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan automation execution: %w", err)
		}
		r.Success = successInt != 0
		records = append(records, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate automation executions: %w", err)
	}
	return records, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

