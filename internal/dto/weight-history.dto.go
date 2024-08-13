package dto

import (
	"errors"
	"time"
)

type CreateWeightHistoryRequest struct {
	Weight float64 `json:"weight" validate:"required"`
	Date   string  `json:"date"`
}

type WeightHistoryResponse struct {
	Weight float64 `json:"weight"`
	Date   string  `json:"date"`
}

type FilterGetWeightHistory struct {
	DateFrom  string
	DateTo    string
	Page      int
	Limit     int
	IsCurrent bool
}

func (r CreateWeightHistoryRequest) Validate() error {
	now := time.Now()
	date, _ := time.Parse("2006-01-02", r.Date)
	dateNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weightDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	if weightDate.After(dateNow) {
		return errors.New("record weight date should not be in the future")
	}

	return nil
}
