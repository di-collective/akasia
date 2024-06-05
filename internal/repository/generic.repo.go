package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"monorepo/pkg/common"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

func NewRepository[Model, ID any](db *sqlx.DB, tableName string) *Repository[Model, ID] {
	return &Repository[Model, ID]{
		db:        db,
		tableName: tableName,
	}
}

type Repository[Model, ID any] struct {
	common.Repository[Model, ID]
	db        *sqlx.DB
	tableName string
}

// Get retrieves a single entity by its ID.
func (repo *Repository[Model, ID]) Get(ctx context.Context, id ID, txs ...*sql.Tx) (*Model, error) {
	stmt, args, err := goqu.Dialect("postgres").
		From(repo.tableName).
		Where(
			goqu.C("id").Eq(id),
			goqu.C("deleted_at").IsNull(),
		).
		Limit(1).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrPreparingStatement, err)
	}

	row := repo.db.QueryRowxContext(ctx, stmt, args...)
	if row.Err() != nil && errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, ErrNoResult
	} else if row.Err() != nil {
		return nil, fmt.Errorf("%w; %w", ErrExecutingStatement, row.Err())
	}

	var entity = new(Model)
	if err := row.StructScan(entity); err != nil {
		return nil, fmt.Errorf("%w; %w", ErrScanResult, err)
	}

	return entity, nil
}

// List retrieves a list of entities based on filters and options.
func (repo *Repository[Model, ID]) List(ctx context.Context, opt *common.FilterOptions, txs ...*sql.Tx) ([]*Model, error) {
	var filter []exp.Expression
	if opt == nil || opt.Filter == nil || len(opt.Filter) <= 0 {
		filter = []exp.Expression{goqu.C("deleted_at").IsNull()}
	} else {
		filter = append(opt.Filter, goqu.C("deleted_at").IsNull())
	}

	page := 1
	if opt != nil && opt.Page > 0 {
		page = opt.Page
	}

	limit := uint(10)
	if opt != nil && opt.Limit > 0 {
		limit = uint(opt.Limit)
	}

	offset := uint(page-1) * limit

	stmt, args, err := goqu.Dialect("postgres").
		From(repo.tableName).
		Where(filter...).
		Select(opt.Select...).
		Order(opt.Sort...).
		Offset(offset).
		Limit(limit).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrPreparingStatement, err)
	}

	rows, err := repo.db.QueryxContext(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrExecutingStatement, err)
	} else if rows.Err() != nil && errors.Is(rows.Err(), sql.ErrNoRows) {
		return nil, ErrNoResult
	} else if rows.Err() != nil {
		return nil, fmt.Errorf("%w; %w", ErrExecutingStatement, rows.Err())
	}
	defer rows.Close()

	var entities = make([]*Model, 0, limit)
	for rows.Next() {
		var entity = new(Model)
		if err := rows.StructScan(entity); err != nil {
			return nil, fmt.Errorf("%w; %w", ErrScanResult, err)
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// Create creates a new entity.
func (repo *Repository[Model, ID]) Create(ctx context.Context, entity *Model, txs ...*sql.Tx) error {
	stmt, args, err := goqu.Dialect("postgres").
		Insert(repo.tableName).
		Rows(entity).
		ToSQL()

	if err != nil {
		return fmt.Errorf("%w; %w", ErrPreparingStatement, err)
	}

	var tx *sql.Tx
	if len(txs) > 0 {
		tx = txs[0]
	} else {
		tx = repo.db.MustBegin().Tx
	}

	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("%w; %w", ErrExecutingStatement, err)
	}

	tx.Commit()
	return nil
}

// Update updates an existing entity.
func (repo *Repository[Model, ID]) Update(ctx context.Context, id ID, entity *Model, txs ...*sql.Tx) error {
	stmt, args, err := goqu.Dialect("postgres").
		Update(repo.tableName).
		Set(entity).
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("%w; %w", ErrPreparingStatement, err)
	}

	var tx *sql.Tx
	if len(txs) > 0 {
		tx = txs[0]
	} else {
		tx = repo.db.MustBegin().Tx
	}

	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("%w; %w", ErrExecutingStatement, err)
	}

	tx.Commit()
	return nil
}

// Delete deletes an entity by its ID.
func (repo *Repository[Model, ID]) Delete(ctx context.Context, id ID, txs ...*sql.Tx) error {
	stmt, args, err := goqu.Dialect("postgres").
		Update(repo.tableName).
		Set(goqu.Record{"deleted_at": time.Now()}).
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("%w; %w", ErrPreparingStatement, err)
	}

	var tx *sql.Tx
	if len(txs) > 0 {
		tx = txs[0]
	} else {
		tx = repo.db.MustBegin().Tx
	}

	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("%w; %w", ErrExecutingStatement, err)
	}

	tx.Commit()
	return nil
}

func (repo *Repository[Model, ID]) Raw(ctx context.Context, statement string, args ...any) (sql.Result, error) {
	tx := repo.db.MustBegin().Tx
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, statement)
	tx.Commit()

	return res, err
}
