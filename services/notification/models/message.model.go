package models

import (
	"database/sql"
	"time"

	"github.com/orsinium-labs/enum"
)

type MessageType enum.Member[string]

var (
	MessageTypes = enum.New(MessageType{"all"}, MessageType{"group"}, MessageType{"individual"})
)

type Message struct {
	ID          string         `db:"id" goqu:"omitempty"`
	Type        string         `db:"type" goqu:"omitempty"`
	CreatedAt   time.Time      `db:"created_at" goqu:"omitempty"`
	DeletedAt   sql.NullTime   `db:"deleted_at" goqu:"omitempty"`
	ScheduledAt sql.NullTime   `db:"scheduled_at" goqu:"omitempty"`
	SentAt      sql.NullTime   `db:"sent_at" goqu:"omitempty"`
	Criteria    sql.NullString `db:"criteria" goqu:"omitempty"`
	Content     string         `db:"content" goqu:"omitempty"`
}

type UserMessage struct {
	ID        string       `db:"id" goqu:"omitempty"`
	UserID    string       `db:"user_id" goqu:"omitempty"`
	MessageID string       `db:"message_id" goqu:"omitempty"`
	CreatedAt time.Time    `db:"created_at" goqu:"omitempty"`
	DeletedAt sql.NullTime `db:"deleted_at" goqu:"omitempty"`
	ReadAt    sql.NullTime `db:"read_at" goqu:"omitempty"`
}

type ViewUserMessage struct {
	ID      string       `db:"message_id"`
	Content string       `db:"message_content"`
	ReadAt  sql.NullTime `db:"read_at"`
}
