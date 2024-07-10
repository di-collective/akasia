package models

import (
	"database/sql"
	"time"
)

type Location struct {
	ID          string       `db:"id" goqu:"omitempty"`
	ClinicID    string       `db:"clinic_id" goqu:"omitempty"`
	Name        string       `db:"name" goqu:"omitempty"`
	Address     string       `db:"address" goqu:"omitempty"`
	Phone       string       `db:"phone" goqu:"omitempty"`
	OpeningTime string       `db:"opening_time" goqu:"omitempty"`
	ClosingTime string       `db:"closing_time" goqu:"omitempty"`
	Capacity    int32        `db:"capacity" goqu:"omitempty"`
	CreatedAt   time.Time    `db:"created_at" goqu:"omitempty"`
	DeletedAt   sql.NullTime `db:"deleted_at" goqu:"omitempty"`
}
