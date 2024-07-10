package service

import (
	"errors"
	"monorepo/internal/repository"
)

var (
	ErrRepositoryQueryFail  = errors.New("failed to fetch data from repository")
	ErrRepositoryMutateFail = errors.New("failed to mutate data to repository")
	ErrValidationFailed     = errors.New("validation failed")
	ErrClinicExist          = errors.New("data clinic is already exists")
	ErrClinicNotFound       = errors.New("data clinic not found")
	ErrLocationExist        = errors.New("data location is already exists")
	ErrNoResult             = repository.ErrNoResult
)
