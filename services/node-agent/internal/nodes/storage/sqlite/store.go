package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	nodesrepo "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/nodes"
)

type Store struct {
	db *sql.DB
}

func OpenNodes(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("nodes storage path cannot be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open nodes database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS nodes (
		id         TEXT    NOT NULL PRIMARY KEY,
		name       TEXT    NOT NULL,
		hostname   TEXT    NOT NULL,
		address    TEXT    NOT NULL,
		version    TEXT    NOT NULL,
		platform   TEXT    NOT NULL,
		status     TEXT    NOT NULL,
		last_seen  DATETIME NOT NULL,
		created_at DATETIME NOT NULL
	)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create nodes table: %w", err)
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_nodes_last_seen ON nodes(last_seen)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create nodes last_seen index: %w", err)
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create nodes status index: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Save(ctx context.Context, node nodesrepo.Node) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("nodes store is not open")
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO nodes (id, name, hostname, address, version, platform, status, last_seen, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		node.ID, node.Name, node.Hostname, node.Address, node.Version, node.Platform, string(node.Status), node.LastSeen, node.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert node: %w", err)
	}

	return nil
}

func (s *Store) UpdateLastSeen(ctx context.Context, id string, lastSeen time.Time) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("nodes store is not open")
	}

	_, err := s.db.ExecContext(ctx,
		`UPDATE nodes SET last_seen = ? WHERE id = ?`,
		lastSeen, id,
	)
	if err != nil {
		return fmt.Errorf("update node last_seen: %w", err)
	}

	return nil
}

func (s *Store) UpdateStatus(ctx context.Context, id string, status nodesrepo.Status) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("nodes store is not open")
	}

	_, err := s.db.ExecContext(ctx,
		`UPDATE nodes SET status = ? WHERE id = ?`,
		string(status), id,
	)
	if err != nil {
		return fmt.Errorf("update node status: %w", err)
	}

	return nil
}

func (s *Store) Get(ctx context.Context, id string) (nodesrepo.Node, error) {
	if s == nil || s.db == nil {
		return nodesrepo.Node{}, fmt.Errorf("nodes store is not open")
	}

	var node nodesrepo.Node
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, hostname, address, version, platform, status, last_seen, created_at FROM nodes WHERE id = ?`,
		id,
	).Scan(&node.ID, &node.Name, &node.Hostname, &node.Address, &node.Version, &node.Platform, &node.Status, &node.LastSeen, &node.CreatedAt)
	if err != nil {
		return nodesrepo.Node{}, fmt.Errorf("get node: %w", err)
	}

	return node, nil
}

func (s *Store) List(ctx context.Context) ([]nodesrepo.Node, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("nodes store is not open")
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, hostname, address, version, platform, status, last_seen, created_at FROM nodes ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	defer rows.Close()

	var nodes []nodesrepo.Node

	for rows.Next() {
		var node nodesrepo.Node
		if err := rows.Scan(&node.ID, &node.Name, &node.Hostname, &node.Address, &node.Version, &node.Platform, &node.Status, &node.LastSeen, &node.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan node: %w", err)
		}
		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nodes: %w", err)
	}

	return nodes, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("nodes store is not open")
	}

	_, err := s.db.ExecContext(ctx, `DELETE FROM nodes WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("remove node: %w", err)
	}

	return nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
