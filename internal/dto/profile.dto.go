package dto

import (
	"errors"
	"regexp"
	"time"
)

type RequestCreateProfile struct {
	UserID      string `json:"user_id,omitempty"`
	Name        string `json:"name" validate:"required"`
	CountryCode string `json:"country_code" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	NIK         string `json:"nik,omitempty"`
	PhotoUrl    string `json:"photo_url,omitempty"`
}

type ResponseCreateProfile struct {
	ID          string `json:"id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	MedicalID   string `json:"medical_id,omitempty"`
	Name        string `json:"name,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Phone       string `json:"phone,omitempty"`
	NIK         string `json:"nik,omitempty"`
	PhotoUrl    string `json:"photo_url,omitempty"`
}

type ResponseGetProfile struct {
	ID            string    `json:"id,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
	Role          string    `json:"role,omitempty"`
	MedicalID     string    `json:"medical_id,omitempty"`
	Name          string    `json:"name,omitempty"`
	CountryCode   string    `json:"country_code,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	NIK           string    `json:"nik,omitempty"`
	Age           string    `json:"age,omitempty"`
	DOB           time.Time `json:"dob,omitempty"`
	Sex           string    `json:"sex,omitempty"`
	BloodType     string    `json:"blood_type,omitempty"`
	Weight        float64   `json:"weight,omitempty"`
	Height        float64   `json:"height,omitempty"`
	ActivityLevel string    `json:"activity_level,omitempty"`
	Allergies     string    `json:"allergies,omitempty"`
	ECRelation    string    `json:"ec_relation,omitempty"`
	ECName        string    `json:"ec_name,omitempty"`
	ECCountryCode string    `json:"ec_country_code,omitempty"`
	ECPhone       string    `json:"ec_phone,omitempty"`
	PhotoUrl      string    `json:"photo_url,omitempty"`
}

type RequestUpdateProfile struct {
	Age           string    `json:"age,omitempty"`
	DOB           time.Time `json:"dob,omitempty"`
	Sex           string    `json:"sex,omitempty" validate:"omitempty,oneof=Male Female"`
	BloodType     string    `json:"blood_type,omitempty"`
	Weight        float64   `json:"weight,omitempty"`
	Height        float64   `json:"height,omitempty"`
	ActivityLevel string    `json:"activity_level,omitempty"`
	Allergies     string    `json:"allergies,omitempty"`
	ECRelation    string    `json:"ec_relation,omitempty" validate:"omitempty,oneof=Wife Husband Mother Father Sister Brother Aunt Uncle Grandmother Grandfather Cousin Friend Spouse Child Other"`
	ECName        string    `json:"ec_name,omitempty"`
	ECCountryCode string    `json:"ec_country_code,omitempty"`
	ECPhone       string    `json:"ec_phone,omitempty"`
	PhotoUrl      string    `json:"photo_url,omitempty"`
}

func (r RequestCreateProfile) Validate() error {
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

	// validate nik (optional)
	if r.NIK != "" {
		regex, err = regexp.Compile(`^\d{6}([04][1-9]|[1256][0-9]|[37][01])(0[1-9]|1[0-2])\d{2}\d{4}$`)
		if err != nil {
			return err
		}

		nikMatch := regex.MatchString(r.NIK)
		if !nikMatch {
			err = errors.New("nik has an invalid format")
			return err
		}
	}

	return nil
}
