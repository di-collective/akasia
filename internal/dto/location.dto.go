package dto

import (
	"errors"
	"regexp"
	"time"
)

type RequestCreateLocation struct {
	ClinicID    string `json:"clinic_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	OpeningTime string `json:"opening_time" validate:"required"`
	ClosingTime string `json:"closing_time" validate:"required"`
	Capacity    int32  `json:"capacity" validate:"required"`
}

type ResponseCreateLocation struct {
	ID          string `json:"id,omitempty"`
	ClinicID    string `json:"clinic_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Address     string `json:"address,omitempty"`
	Phone       string `json:"phone,omitempty"`
	OpeningTime string `json:"opening_time,omitempty"`
	ClosingTime string `json:"closing_time,omitempty"`
	Capacity    int32  `json:"capacity,omitempty"`
}

type RequestUpdateLocation struct {
	ClinicID    string `json:"clinic_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	OpeningTime string `json:"opening_time" validate:"required"`
	ClosingTime string `json:"closing_time" validate:"required"`
	Capacity    int32  `json:"capacity" validate:"required"`
}

type ResponseUpdateLocation struct {
	ID          string `json:"id,omitempty"`
	ClinicID    string `json:"clinic_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Address     string `json:"address,omitempty"`
	Phone       string `json:"phone,omitempty"`
	OpeningTime string `json:"opening_time,omitempty"`
	ClosingTime string `json:"closing_time,omitempty"`
	Capacity    int32  `json:"capacity,omitempty"`
}

type ResponseGetLocation struct {
	ID          string     `json:"id,omitempty"`
	ClinicID    string     `json:"clinic_id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Address     string     `json:"address,omitempty"`
	Phone       string     `json:"phone,omitempty"`
	OpeningTime string     `json:"opening_time,omitempty"`
	ClosingTime string     `json:"closing_time,omitempty"`
	Capacity    int32      `json:"capacity,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func (r RequestCreateLocation) Validate() error {
	// validate phone
	regex, err := regexp.Compile(`^[1-9]{9,12}$`)
	if err != nil {
		return err
	}

	phoneMatch := regex.MatchString(r.Phone)
	if !phoneMatch {
		err = errors.New("phone has an invalid format")
		return err
	}

	return nil
}

func (r RequestUpdateLocation) Validate() error {
	// validate phone
	regex, err := regexp.Compile(`^[1-9]{9,12}$`)
	if err != nil {
		return err
	}

	phoneMatch := regex.MatchString(r.Phone)
	if !phoneMatch {
		err = errors.New("phone has an invalid format")
		return err
	}

	return nil
}
