package service

import (
	"errors"
	"monorepo/internal/repository"
)

var (
	ErrRepositoryQueryFail  = errors.New("failed to fetch data from repository")
	ErrRepositoryMutateFail = errors.New("failed to mutate data to repository")
	ErrValidationFailed     = errors.New("validation failed")
	ErrNoResult             = repository.ErrNoResult
)
