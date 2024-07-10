package models

import (
	"database/sql"
	"time"
)

type Clinic struct {
	ID        string       `db:"id" goqu:"omitempty"`
	Name      string       `db:"name" goqu:"omitempty"`
	Address   string       `db:"address" goqu:"omitempty"`
	Phone     string       `db:"phone" goqu:"omitempty"`
	Logo      string       `db:"logo" goqu:"omitempty"`
	CreatedAt time.Time    `db:"created_at" goqu:"omitempty"`
	DeletedAt sql.NullTime `db:"deleted_at" goqu:"omitempty"`
}
