package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	notifrepo "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

type Store struct {
	db *sql.DB
}

func OpenNotifications(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("notifications storage path cannot be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open notifications database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS notifications (
		id         TEXT    NOT NULL,
		title      TEXT    NOT NULL,
		message    TEXT    NOT NULL,
		severity   TEXT    NOT NULL,
		source     TEXT    NOT NULL,
		provider   TEXT,
		status     TEXT    NOT NULL,
		created_at DATETIME NOT NULL,
		sent_at    DATETIME,
		error      TEXT
	)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create notifications table: %w", err)
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create notifications created_at index: %w", err)
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create notifications status index: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Save(ctx context.Context, notification notifrepo.Notification) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("notifications store is not open")
	}

	var provider, errorMsg interface{}
	if notification.Provider != "" {
		provider = notification.Provider
	}
	if notification.Error != "" {
		errorMsg = notification.Error
	}

	var sentAt interface{}
	if notification.SentAt != nil {
		sentAt = *notification.SentAt
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO notifications (id, title, message, severity, source, provider, status, created_at, sent_at, error) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		notification.ID,
		notification.Title,
		notification.Message,
		string(notification.Severity),
		string(notification.Source),
		provider,
		string(notification.Status),
		notification.CreatedAt,
		sentAt,
		errorMsg,
	)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}

	return nil
}

func (s *Store) Recent(limit int) ([]notifrepo.Notification, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("notifications store is not open")
	}

	if limit <= 0 {
		limit = 100
	}

	rows, err := s.db.Query(
		`SELECT id, title, message, severity, source, provider, status, created_at, sent_at, error FROM notifications ORDER BY created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query recent notifications: %w", err)
	}
	defer rows.Close()

	var notifications []notifrepo.Notification

	for rows.Next() {
		var n notifrepo.Notification
		var provider, errorMsg sql.NullString
		var sentAt sql.NullTime

		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Message,
			&n.Severity,
			&n.Source,
			&provider,
			&n.Status,
			&n.CreatedAt,
			&sentAt,
			&errorMsg,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}

		if provider.Valid {
			n.Provider = provider.String
		}
		if errorMsg.Valid {
			n.Error = errorMsg.String
		}
		if sentAt.Valid {
			t := sentAt.Time
			n.SentAt = &t
		}

		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notifications: %w", err)
	}

	return notifications, nil
}

func (s *Store) UpdateStatus(ctx context.Context, id string, status notifrepo.Status, sentAt *time.Time, errorMsg string) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("notifications store is not open")
	}

	var sentAtVal interface{}
	if sentAt != nil {
		sentAtVal = *sentAt
	}

	_, err := s.db.ExecContext(ctx,
		`UPDATE notifications SET status = ?, sent_at = ?, error = ? WHERE id = ?`,
		string(status),
		sentAtVal,
		errorMsg,
		id,
	)
	if err != nil {
		return fmt.Errorf("update notification status: %w", err)
	}

	return nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
