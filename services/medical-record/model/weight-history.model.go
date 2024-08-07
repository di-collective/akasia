package model

import "time"

type WeightHistory struct {
	ProfileID string    `db:"profile_id" goqu:"omitempty"`
	Weight    float64   `db:"weight" goqu:"omitempty"`
	CreatedAt time.Time `db:"created_at" goqu:"omitempty"`
}
