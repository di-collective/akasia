package dto

import (
	"errors"
	"time"
)

type CreateWeightGoalRequest struct {
	StartingWeight float64 `json:"starting_weight,omitempty" validate:"required"`
	TargetWeight   float64 `json:"target_weight,omitempty" validate:"required"`
	ActivityLevel  string  `json:"activity_level,omitempty" validate:"required"`
	Pace           string  `json:"pace,omitempty" validate:"required"`
}

type CreateWeightGoalResponse struct {
	StartingWeight      float64 `json:"starting_weight,omitempty"`
	StartingDate        string  `json:"starting_date,omitempty"`
	TargetWeight        float64 `json:"target_weight,omitempty"`
	TargetDate          string  `json:"target_date,omitempty"`
	ActivityLevel       string  `json:"activity_level,omitempty"`
	DailyCaloriesBudget float64 `json:"daily_calories_budget,omitempty"`
	CaloriesToMaintain  float64 `json:"calories_to_maintain,omitempty"`
	Flag                string  `json:"flag,omitempty"`
	Pace                string  `json:"pace,omitempty"`
}

type GetWeightGoalResponse struct {
	StartingWeight      float64 `json:"starting_weight,omitempty"`
	StartingDate        string  `json:"starting_date,omitempty"`
	TargetWeight        float64 `json:"target_weight,omitempty"`
	TargetDate          string  `json:"target_date,omitempty"`
	ActivityLevel       string  `json:"activity_level,omitempty"`
	DailyCaloriesBudget float64 `json:"daily_calories_budget,omitempty"`
	CaloriesToMaintain  float64 `json:"calories_to_maintain,omitempty"`
	Flag                string  `json:"flag,omitempty"`
	Pace                string  `json:"pace,omitempty"`
}

type UpdateWeightGoalRequest struct {
	CurrentWeight  float64 `json:"current_weight,omitempty"`
	StartingWeight float64 `json:"starting_weight,omitempty"`
	StartingDate   string  `json:"starting_date,omitempty"`
	TargetWeight   float64 `json:"target_weight,omitempty"`
	ActivityLevel  string  `json:"activity_level,omitempty"`
	Pace           string  `json:"pace,omitempty"`
}

type WeightGoalPace struct {
	Pace                string  `json:"pace,omitempty"`
	DailyCaloriesBudget float64 `json:"daily_calories_budget,omitempty"`
	TargetDate          string  `json:"target_date,omitempty"`
}
type SimulationWeightGoalRequest struct {
	StartingWeight float64 `json:"starting_weight,omitempty" validate:"required"`
	TargetWeight   float64 `json:"target_weight,omitempty" validate:"required"`
	ActivityLevel  string  `json:"activity_level,omitempty" validate:"required"`
}

type SimulationWeightGoalResponse struct {
	StartingWeight     float64          `json:"starting_weight,omitempty"`
	StartingDate       string           `json:"starting_date,omitempty"`
	TargetWeight       float64          `json:"target_weight,omitempty"`
	ActivityLevel      string           `json:"activity_level,omitempty"`
	CaloriesToMaintain float64          `json:"calories_to_maintain,omitempty"`
	Flag               string           `json:"flag,omitempty"`
	Pacing             []WeightGoalPace `json:"pacing,omitempty"`
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

func (r UpdateWeightGoalRequest) Validate() error {
	now := time.Now()
	startTime, _ := time.Parse("2006-01-02", r.StartingDate)

	dateNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())

	if startDate.After(dateNow) {
		return errors.New("start date should be today or in the past")
	}

	return nil
}

func (r SimulationWeightGoalRequest) Validate() error {
	if r.TargetWeight == r.StartingWeight {
		return errors.New("target weight must be different from starting weight")
	}

	if r.TargetWeight <= 0 || r.StartingWeight <= 0 {
		return errors.New("target weight must be gtreather than 0")
	}

	return nil
}
