package service

import (
	"errors"
	"monorepo/internal/repository"
)

var (
	ErrRepositoryQueryFail   = errors.New("failed to fetch data from repository")
	ErrRepositoryMutateFail  = errors.New("failed to mutate data to repository")
	ErrValidationFailed      = errors.New("validation failed")
	ErrPasswordHashingFailed = errors.New("failed to hash password")
	ErrNoResult              = repository.ErrNoResult
)
