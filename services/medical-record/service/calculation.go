package service

import (
	"math"
	"monorepo/internal/constants"
	"strings"
	"time"
)

func CalculateCalorieBudget(gender string, weight, height float64, age int, activityLevel string) float64 {
	var (
		bmr            float64
		activityFactor float64
	)

	if strings.ToLower(gender) == constants.GenderMale {
		bmr = 66.5 + (13.75 * weight) + (5.003 * height) - (6.75 * float64(age))
	} else if strings.ToLower(gender) == constants.GenderFemale {
		bmr = 655.1 + (9.563 * weight) + (1.850 * height) - (4.676 * float64(age))
	}

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

func CalculateWeeksToTarget(startWeight, targetWeight, weightChangePerWeek float64, flag string) float64 {
	if strings.ToLower(flag) == constants.WeightGoalLoss {
		return math.Ceil((startWeight - targetWeight) / weightChangePerWeek)
	} else if strings.ToLower(flag) == constants.WeightGoalGain {
		return math.Ceil((targetWeight - startWeight) / weightChangePerWeek)
	}

	return 0.0 // if flag not valid
}

func CalculateTargetDate(gender string, age int, height, startWeight, targetWeight, weightChangePerWeek float64, activityLevel, flag string) time.Time {
	weeksToTarget := CalculateWeeksToTarget(startWeight, targetWeight, weightChangePerWeek, flag)
	now := time.Now()
	targetDate := now.AddDate(0, 0, int(weeksToTarget*7))
	return targetDate
}
