package model

import (
	"database/sql"
	"time"
)

type WeightHistory struct {
	ID        string       `db:"id" goqu:"omitempty" json:"id"`
	ProfileID string       `db:"profile_id" goqu:"omitempty"`
	Weight    float64      `db:"weight" goqu:"omitempty"`
	CreatedAt time.Time    `db:"created_at" goqu:"omitempty"`
	UpdatedAt sql.NullTime `db:"updated_at" goqu:"omitempty"`
	DeletedAt sql.NullTime `db:"deleted_at" goqu:"omitempty" json:"deleted_at"`
}
