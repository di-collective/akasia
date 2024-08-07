package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrPreparingStatement   = errors.New("failed to prepare SQL statement")
	ErrExecutingStatement   = errors.New("failed to execute SQL statement")
	ErrScanResult           = errors.New("failed to scan SQL result")
	ErrNoResult             = sql.ErrNoRows
	ErrRepositoryQueryFail  = errors.New("failed to fetch data from repository")
	ErrRepositoryMutateFail = errors.New("failed to mutate data to repository")
	ErrValidationFailed     = errors.New("validation failed")
	ErrExist                = errors.New("data is already exists")
)
