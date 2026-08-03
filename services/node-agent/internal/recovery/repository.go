package recovery

import "context"

type Repository interface {
	Save(ctx context.Context, execution Execution) error
	Recent(limit int) ([]Execution, error)
	Close() error
}
