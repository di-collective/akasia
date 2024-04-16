package dto

import (
	"time"
)

type Object[T any] struct {
	Data    *T     `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type NotificationMessage struct {
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	SentAt      *time.Time `json:"sent_at"`
	Type        string     `json:"type"`
	Criteria    string     `json:"criteria"`
	Content     string     `json:"content"`
}

type NotificationUserMessage struct {
	ID      string     `json:"id"`
	Content string     `json:"content"`
	ReadAt  *time.Time `json:"read_at"`
}

type RequestListNotificationMessage struct {
	IDs     []string `schema:"ids"`
	Type    []string `schema:"type"`
	Content string   `schema:"content"`
}

type RequestMutateNotificationMessage struct {
	ScheduledAt *time.Time `json:"scheduled_at" validate:"omitempty"`
	Type        string     `json:"type" validate:"required,oneof=all group individual"`
	Criteria    string     `json:"criteria" validate:"omitempty"`
	Content     string     `json:"content" validate:"required"`
}
