package dto

import (
	"errors"
	"monorepo/internal/constants"
	"time"
)

type RequestCreateEvent struct {
	LocationID string    `json:"location_id" validate:"required"`
	Status     string    `json:"status,omitempty"`
	Type       string    `json:"type" validate:"required,oneof=holiday appointment"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time"`
}

type ResponseCreateEvent struct {
	ID         string    `json:"id,omitempty"`
	ProfileID  *string   `json:"profile_id,omitempty"`
	LocationID string    `json:"location_id,omitempty"`
	Status     string    `json:"status,omitempty"`
	Type       string    `json:"type,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	EndTime    time.Time `json:"end_time,omitempty"`
}

type ResponseDetailEvent struct {
	ID         string    `json:"id,omitempty"`
	ProfileID  *string   `json:"profile_id,omitempty"`
	LocationID string    `json:"location_id,omitempty"`
	Status     string    `json:"status,omitempty"`
	Type       string    `json:"type,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	EndTime    time.Time `json:"end_time,omitempty"`
	Capacity   int       `json:"capacity,omitempty"`
}

type ResponseGetEvents struct {
	Capacity int                   `json:"capacity,omitempty"`
	Events   []ResponseDetailEvent `json:"events,omitempty"`
}

type FilterGetEvents struct {
	Page       int
	Limit      int
	LocationID string
	StartTime  time.Time
	EndTime    time.Time
}

func (r RequestCreateEvent) Validate() error {
	if r.Type != constants.Appointment && r.Status != "" {
		err := errors.New("status should only be provided when type is appointment")
		return err
	}

	return nil
}
