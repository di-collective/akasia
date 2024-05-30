package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string       `db:"id" goqu:"omitempty"`
	Provider  string       `db:"provider" goqu:"omitempty"`
	Handle    string       `db:"handle" goqu:"omitempty"`
	Password  string       `db:"password" goqu:"omitempty"`
	CreatedAt time.Time    `db:"created_at" goqu:"omitempty"`
	DeletedAt sql.NullTime `db:"deleted_at" goqu:"omitempty"`
}

type Profile struct {
	ID          string         `db:"id" json:"id" goqu:"omitempty"`
	UserID      string         `db:"user_id" json:"user_id" goqu:"omitempty"`
	MedicalID   string         `db:"medical_id" json:"medical_id" goqu:"omitempty"`
	FirstName   string         `db:"first_name" json:"first_name" goqu:"omitempty"`
	LastName    string         `db:"last_name" json:"last_name" goqu:"omitempty"`
	CountryCode string         `db:"country_code" json:"country_code" goqu:"omitempty"`
	Phone       string         `db:"phone" json:"phone" goqu:"omitempty"`
	NIK         sql.NullString `db:"nik" json:"nik" goqu:"omitempty"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at" goqu:"omitempty"`
	DeletedAt   sql.NullTime   `db:"deleted_at" json:"deleted_at" goqu:"omitempty"`
}
