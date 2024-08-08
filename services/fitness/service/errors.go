package service

import (
	"errors"
	"monorepo/internal/repository"
)

var (
	ErrRepositoryQueryFail  = errors.New("failed to fetch data from repository")
	ErrRepositoryMutateFail = errors.New("failed to mutate data to repository")
	ErrNoResult             = repository.ErrNoResult
	ErrGetProfile           = errors.New("failed to get profile")
	ErrUpdateProfile        = errors.New("failed to update profile")
	ErrWeightGoalExist      = errors.New("data weight goal is already exists")
	ErrNotFound             = errors.New("data not found")
)
