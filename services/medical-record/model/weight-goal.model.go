package model

import (
	"database/sql"
	"time"
)

type WeightGoal struct {
	ID             string       `db:"id" goqu:"omitempty"`
	ProfileID      string       `db:"profile_id" goqu:"omitempty"`
	StartingWeight float64      `db:"starting_weight" goqu:"omitempty"`
	StartingDate   time.Time    `db:"starting_date" goqu:"omitempty"`
	TargetWeight   float64      `db:"target_weight" goqu:"omitempty"`
	TargetDate     time.Time    `db:"target_date" goqu:"omitempty"`
	CalorieBudget  float64      `db:"calorie_budget" goqu:"omitempty"`
	Flag           string       `db:"flag" goqu:"omitempty"` // gain | loss | maintain
	ActivityLevel  string       `db:"activity_level" goqu:"omitempty"`
	CreatedAt      time.Time    `db:"created_at" goqu:"omitempty"`
	UpdatedAt      sql.NullTime `db:"updated_at" goqu:"omitempty"`
	DeletedAt      sql.NullTime `db:"deleted_at" goqu:"omitempty"`
}
