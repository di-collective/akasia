package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrPreparingStatement = errors.New("failed to prepare SQL statement")
	ErrExecutingStatement = errors.New("failed to execute SQL statement")
	ErrScanResult         = errors.New("failed to scan SQL result")
	ErrNoResult           = sql.ErrNoRows
)
