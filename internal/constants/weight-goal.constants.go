package constants

const (
	WeightGoalGain     = "gain"
	WeightGoalLoss     = "loss"
	WeightGoalMaintain = "maintain"

	GenderMale   = "male"
	GenderFemale = "female"
)

const (
	ActivityLevelSedentary      = "Sedentary"
	ActivityLevelLightActive    = "Lightly Active"
	ActivityLevelModerateActive = "Moderately Active"
	ActivityLevelVeryActive     = "Very Active"

	ActivityLevelSedentaryVal      float64 = 1.2
	ActivityLevelLightActiveVal    float64 = 1.375
	ActivityLevelModerateActiveVal float64 = 1.55
	ActivityLevelVeryActiveVal     float64 = 1.725
)

const (
	WeeklyWeightPaceRelaxed = "relaxed"
	WeeklyWeightPaceNormal  = "normal"
	WeeklyWeightStrict      = "strict"

	WeeklyWeightPaceRelaxedVal float64 = 0.25
	WeeklyWeightPaceNormalVal  float64 = 0.5
	WeeklyWeightPaceStrictVal  float64 = 1
)
