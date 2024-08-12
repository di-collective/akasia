package model

import (
	"database/sql"
	"time"
)

type WeightGoal struct {
	ID                 string       `db:"id" goqu:"omitempty" json:"id"`
	ProfileID          string       `db:"profile_id" goqu:"omitempty" json:"profile_id"`
	StartingWeight     float64      `db:"starting_weight" goqu:"omitempty" json:"starting_weight"`
	StartingDate       time.Time    `db:"starting_date" goqu:"omitempty" json:"starting_date"`
	TargetWeight       float64      `db:"target_weight" goqu:"omitempty" json:"target_weight"`
	TargetDate         time.Time    `db:"target_date" goqu:"omitempty" json:"target_date"`
	DailyCalorieBudget float64      `db:"daily_calories_budget" goqu:"omitempty" json:"daily_calories_budget"`
	CaloriesToMaintain float64      `db:"calories_to_maintain" goqu:"omitempty" json:"calories_to_maintain"`
	Flag               string       `db:"flag" goqu:"omitempty" json:"flag"` // gain | loss | maintain
	ActivityLevel      string       `db:"activity_level" goqu:"omitempty" json:"activity_level"`
	Pace               string       `db:"pace" goqu:"omitempty" json:"pace"`
	CreatedAt          time.Time    `db:"created_at" goqu:"omitempty" json:"created_at"`
	UpdatedAt          sql.NullTime `db:"updated_at" goqu:"omitempty" json:"updated_at"`
	DeletedAt          sql.NullTime `db:"deleted_at" goqu:"omitempty" json:"deleted_at"`
}
