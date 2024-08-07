package dto

import (
	"errors"
)

type CreateWeightGoalRequest struct {
	StartingWeight float64 `json:"starting_weight,omitempty"`
	TargetWeight   float64 `json:"target_weight,omitempty"`
	ActivityLevel  string  `json:"activity_level,omitempty"`
}

type CreateWeightGoalResponse struct {
	TargetWeight float64 `json:"target_weight,omitempty"`
	TargetDate   string  `json:"target_date,omitempty"`
}

type GetWeightGoalResponse struct {
	StartingWeight float64 `json:"starting_weight,omitempty"`
	StartingDate   string  `json:"starting_date,omitempty"`
	TargetWeight   float64 `json:"target_weight,omitempty"`
	TargetDate     string  `json:"target_date,omitempty"`
	ActivityLevel  string  `json:"activity_level,omitempty"`
	CalorieBudget  float64 `json:"calorie_budget,omitempty"`
	Flag           string  `json:"flag,omitempty"`
}

func (r CreateWeightGoalRequest) Validate() error {
	if r.TargetWeight == r.StartingWeight {
		return errors.New("target weight must be different from starting weight")
	}

	if r.TargetWeight <= 0 || r.StartingWeight <= 0 {
		return errors.New("target weight must be gtreather than 0")
	}

	return nil
}
