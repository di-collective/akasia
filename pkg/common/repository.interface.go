package common

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9/exp"
)

type FilterOptions struct {
	Filter []exp.Expression
	Sort   []exp.Expression
	Select []any
	Page   int
	Limit  int
}

type Repository[Model, ID any] interface {
	// Get retrieves a single entity by its ID.
	Get(ctx context.Context, id ID, tx ...*sql.Tx) (*Model, error)

	// List retrieves a list of entities based on filters and options.
	List(ctx context.Context, opt *FilterOptions, tx ...*sql.Tx) ([]*Model, error)

	// Create creates a new entity.
	Create(ctx context.Context, entity *Model, tx ...*sql.Tx) error

	// Update updates an existing entity.
	Update(ctx context.Context, id ID, entity *Model, tx ...*sql.Tx) error

	// Delete deletes an entity by its ID.
	Delete(ctx context.Context, id ID, tx ...*sql.Tx) error

	// Raw execute raw SQL
	Raw(ctx context.Context, statement string, args ...any) (sql.Result, error)
}
