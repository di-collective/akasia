package models

import (
	"database/sql"
	"time"
)

type Event struct {
	ID         string       `db:"id" goqu:"omitempty"`
	ProfileID  *string      `db:"profile_id" goqu:"omitempty"`
	LocationID string       `db:"location_id" goqu:"omitempty"`
	Status     string       `db:"status" goqu:"omitempty"`
	Type       string       `db:"type" goqu:"omitempty"`
	StartTime  time.Time    `db:"start_time" goqu:"omitempty"`
	EndTime    time.Time    `db:"end_time" goqu:"omitempty"`
	CreatedAt  time.Time    `db:"created_at" goqu:"omitempty"`
	DeletedAt  sql.NullTime `db:"deleted_at" goqu:"omitempty"`
}
