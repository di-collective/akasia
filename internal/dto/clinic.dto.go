package dto

import (
	"errors"
	"regexp"
	"time"
)

type RequestCreateClinic struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Logo    string `json:"logo" validate:"required"`
}

type ResponseCreateClinic struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Logo    string `json:"logo,omitempty"`
}

type RequestUpdateClinic struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Logo    string `json:"logo" validate:"required"`
}

type ResponseUpdateClinic struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Logo    string `json:"logo,omitempty"`
}

type ResponseGetClinic struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Address   string     `json:"address,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Logo      string     `json:"logo,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type FilterGetClinic struct {
	Page  int
	Limit int
}

func (r RequestCreateClinic) Validate() error {
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

func (r RequestUpdateClinic) Validate() error {
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
