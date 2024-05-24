package models

import (
	"database/sql"
	"time"
)

type ResetPassword struct {
	ID         string       `db:"id" goqu:"omitempty"`
	UserID     string       `db:"user_id" goqu:"omitempty"`
	ResetToken string       `db:"reset_token" goqu:"omitempty"`
	IsUsed     bool         `db:"is_used" goqu:"omitempty"`
	CreatedAt  time.Time    `db:"created_at" goqu:"omitempty"`
	DeletedAt  sql.NullTime `db:"deleted_at" goqu:"omitempty"`
}
