package service

import (
	"math"
	"monorepo/internal/constants"
	"strings"
	"time"
)

// Calculate BMR (basal metabolic rate) using Harris Benedict formula
func CalculateBMR(gender string, weight, height float64, age int) float64 {
	if strings.ToLower(gender) == constants.GenderMale {
		return 66.5 + (13.75 * weight) + (5.003 * height) - (6.75 * float64(age))
	} else if strings.ToLower(gender) == constants.GenderFemale {
		return 655.1 + (9.563 * weight) + (1.850 * height) - (4.676 * float64(age))
	}

	return 0.0
}

// Calculating daily calorie budget based on activity level
func CalculateCalorieToMaintain(gender string, weight, height float64, age int, activityLevel string) float64 {
	bmr := CalculateBMR(gender, weight, height, age)

	var activityFactor float64
	switch activityLevel {
	case constants.ActivityLevelSedentary:
		activityFactor = constants.ActivityLevelSedentaryVal
	case constants.ActivityLevelLightActive:
		activityFactor = constants.ActivityLevelModerateActiveVal
	case constants.ActivityLevelModerateActive:
		activityFactor = constants.ActivityLevelLightActiveVal
	case constants.ActivityLevelVeryActive:
		activityFactor = constants.ActivityLevelVeryActiveVal
	default:
		activityFactor = constants.ActivityLevelSedentaryVal
	}

	return bmr * activityFactor
}

// Calciulate daily deficit or surplus calorie
func CalculateDailyCalorieBudget(calorieToMaintain, weightChangePerWeek float64, flag string) float64 {
	caloriesPerKg := 7700.0 // 1 kg lemak kira-kira setara dengan 7700 kalori
	dailyCalorieChange := (weightChangePerWeek * caloriesPerKg) / 7

	if strings.ToLower(flag) == constants.WeightGoalLoss {
		return calorieToMaintain - dailyCalorieChange
	} else if strings.ToLower(flag) == constants.WeightGoalGain {
		return calorieToMaintain + dailyCalorieChange
	}
	return calorieToMaintain // if flag not valid
}

// Calculate days target based on goal
func CalculateDaysToTarget(startWeight, targetWeight, weightChangePerWeek float64, flag string) float64 {
	weightChangePerDay := weightChangePerWeek / 7

	if strings.ToLower(flag) == constants.WeightGoalLoss {
		return math.Ceil((startWeight - targetWeight) / weightChangePerDay)
	} else if strings.ToLower(flag) == constants.WeightGoalGain {
		return math.Ceil((targetWeight - startWeight) / weightChangePerDay)
	}

	return 0.0 // if flag not valid
}

// Calculate goal target date
func CalculateTargetDate(gender string, age int, height, startWeight, targetWeight, weightChangePerWeek float64, activityLevel, flag string, startdate time.Time) time.Time {
	daysToTarget := CalculateDaysToTarget(startWeight, targetWeight, weightChangePerWeek, flag)
	targetDate := startdate.AddDate(0, 0, int(daysToTarget))
	return targetDate
}

// Chek if newly record weight achieve goal weight
func IsAchieveGoal(recordedWeight, targetWeight float64, flag string) bool {
	if strings.ToLower(flag) == constants.WeightGoalGain {
		return recordedWeight >= targetWeight
	} else if strings.ToLower(flag) == constants.WeightGoalLoss {
		return recordedWeight <= targetWeight
	}
	return false
}
