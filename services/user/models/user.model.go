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
	ID          string       `db:"id" goqu:"omitempty"`
	UserID      string       `db:"user_id" goqu:"omitempty"`
	MedicalID   string       `db:"medical_id" goqu:"omitempty"`
	FirstName   string       `db:"first_name" goqu:"omitempty"`
	LastName    string       `db:"last_name" goqu:"omitempty"`
	CountryCode string       `db:"country_code" goqu:"omitempty"`
	Phone       string       `db:"phone" goqu:"omitempty"`
	NIK         string       `db:"nik" goqu:"omitempty"`
	CreatedAt   time.Time    `db:"created_at" goqu:"omitempty"`
	DeletedAt   sql.NullTime `db:"deleted_at" goqu:"omitempty"`
}
