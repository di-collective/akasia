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
	ID            string       `db:"id" json:"id" goqu:"omitempty"`
	UserID        string       `db:"user_id" json:"user_id" goqu:"omitempty"`
	MedicalID     string       `db:"medical_id" json:"medical_id" goqu:"omitempty"`
	FirstName     string       `db:"first_name" json:"first_name" goqu:"omitempty"`
	LastName      string       `db:"last_name" json:"last_name" goqu:"omitempty"`
	CountryCode   string       `db:"country_code" json:"country_code" goqu:"omitempty"`
	Phone         string       `db:"phone" json:"phone" goqu:"omitempty"`
	NIK           *string      `db:"nik" json:"nik" goqu:"omitempty"`
	Age           *string      `db:"age" json:"age" goqu:"omitempty"`
	DOB           *time.Time   `db:"dob" json:"dob" goqu:"omitempty"`
	Sex           *string      `db:"sex" json:"sex" goqu:"omitempty"`
	BloodType     *string      `db:"blood_type" json:"blood_type" goqu:"omitempty"`
	Weight        *float64     `db:"weight" json:"weight" goqu:"omitempty"`
	Height        *float64     `db:"height" json:"height" goqu:"omitempty"`
	ActivityLevel *string      `db:"activity_level" json:"activity_level" goqu:"omitempty"`
	Allergies     *string      `db:"allergies" json:"allergies" goqu:"omitempty"`
	ECRelation    *string      `db:"ec_relation" json:"ec_relation" goqu:"omitempty"`
	ECName        *string      `db:"ec_name" json:"ec_name" goqu:"omitempty"`
	ECCountryCode *string      `db:"ec_country_code" json:"ec_country_code" goqu:"omitempty"`
	ECPhone       *string      `db:"ec_phone" json:"ec_phone" goqu:"omitempty"`
	PhotoUrl      *string      `db:"photo_url" json:"photo_url" goqu:"omitempty"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at" goqu:"omitempty"`
	DeletedAt     sql.NullTime `db:"deleted_at" json:"deleted_at" goqu:"omitempty"`
}
