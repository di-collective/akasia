package dto

import (
	"errors"
	"regexp"
)

type RequestCreateProfile struct {
	UserID      string `json:"user_id,omitempty"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	CountryCode string `json:"country_code" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	NIK         string `json:"nik" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	BirthDate   string `json:"dob" validate:"required"`
}

func (r RequestCreateProfile) Validate() error {
	// validate phone
	regex, err := regexp.Compile(`^[1-9].{9,12}$`)
	if err != nil {
		return err
	}

	phoneMatch := regex.MatchString(r.Phone)
	if !phoneMatch {
		err = errors.New("phone has an invalid format")
		return err
	}

	// validate nik
	regex, err = regexp.Compile(`^\d{6}([04][1-9]|[1256][0-9]|[37][01])(0[1-9]|1[0-2])\d{2}\d{4}$`)
	if err != nil {
		return err
	}

	nikMatch := regex.MatchString(r.NIK)
	if !nikMatch {
		err = errors.New("nik has an invalid format")
		return err
	}

	return nil
}
