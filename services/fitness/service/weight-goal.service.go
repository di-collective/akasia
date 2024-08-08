package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"monorepo/internal/constants"
	"monorepo/internal/dto"
	"monorepo/pkg/common"
	"monorepo/services/fitness/model"
	"strconv"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

const shortdDateLayout = "2006-01-02"

func NewWeightGoalService(
	tbWeightGoal common.Repository[model.WeightGoal, string],
	tbWeightHistory common.Repository[model.WeightHistory, string],
) *WeightGoalService {
	service := &WeightGoalService{}
	service.validate = validator.New()
	service.tables.weightGoal = tbWeightGoal
	service.tables.weightHistory = tbWeightHistory

	return service
}

type WeightGoalService struct {
	validate *validator.Validate
	tables   struct {
		weightGoal    common.Repository[model.WeightGoal, string]
		weightHistory common.Repository[model.WeightHistory, string]
	}
}

func (service *WeightGoalService) IsWeightGoalExists(ctx context.Context, profileID string) (bool, error) {
	existing, err := service.tables.weightGoal.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("profile_id").Eq(profileID),
		},
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return false, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	return len(existing) > 0, nil
}

func (service *WeightGoalService) CreateWightGoal(ctx context.Context, body dto.CreateWeightGoalRequest) (*dto.CreateWeightGoalResponse, error) {
	now := time.Now()
	profile, err := GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrGetProfile, err)
	}

	// check exist
	isExist, err := service.IsWeightGoalExists(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	if isExist {
		return nil, fmt.Errorf("%w", ErrWeightGoalExist)
	}

	var wgFlag string
	if body.TargetWeight > body.StartingWeight {
		wgFlag = constants.WeightGoalGain
	} else {
		wgFlag = constants.WeightGoalLoss
	}

	// calculate calorie budget
	age, _ := strconv.Atoi(profile.Age)

	var paceVal float64
	switch strings.ToLower(body.Pace) {
	case constants.WeeklyWeightPaceRelaxed:
		paceVal = constants.WeeklyWeightPaceRelaxedVal
	case constants.WeeklyWeightPaceNormal:
		paceVal = constants.WeeklyWeightPaceNormalVal
	case constants.WeeklyWeightStrict:
		paceVal = constants.WeeklyWeightPaceStrictVal
	}

	caloriesToMaintain := math.Round(CalculateCalorieToMaintain(profile.Sex, body.StartingWeight, profile.Height, age, body.ActivityLevel)*100) / 100
	dailyCaloriesBudget := math.Round(CalculateDailyCalorieBudget(caloriesToMaintain, paceVal, wgFlag)*100) / 100
	targetDate := CalculateTargetDate(profile.Sex, age, profile.Height, body.StartingWeight, body.TargetWeight, paceVal, body.ActivityLevel, wgFlag, now)

	newWeightGoal := &model.WeightGoal{
		ID:                 ulid.Make().String(),
		ProfileID:          profile.ID,
		StartingWeight:     body.StartingWeight,
		StartingDate:       now,
		TargetWeight:       body.TargetWeight,
		TargetDate:         targetDate,
		DailyCalorieBudget: dailyCaloriesBudget,
		CaloriesToMaintain: caloriesToMaintain,
		Flag:               wgFlag,
		CreatedAt:          now,
		ActivityLevel:      body.ActivityLevel,
		Pace:               body.Pace,
	}

	// insert weight goal
	if err := service.tables.weightGoal.Create(ctx, newWeightGoal); err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	newWeightHistory := &model.WeightHistory{
		ProfileID: newWeightGoal.ProfileID,
		Weight:    body.StartingWeight,
		CreatedAt: time.Now(),
	}

	// insert weight history
	if err := service.tables.weightHistory.Create(ctx, newWeightHistory); err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	// update profile
	if err := UpdateProfile(ctx, profile.UserID, dto.RequestUpdateProfile{
		Weight:        body.StartingWeight,
		ActivityLevel: body.ActivityLevel,
	}); err != nil {
		return nil, fmt.Errorf("%w; %w", ErrUpdateProfile, err)
	}

	return &dto.CreateWeightGoalResponse{
		StartingWeight:      newWeightGoal.StartingWeight,
		StartingDate:        newWeightGoal.StartingDate.Format(shortdDateLayout),
		TargetWeight:        newWeightGoal.TargetWeight,
		TargetDate:          newWeightGoal.TargetDate.Format(shortdDateLayout),
		ActivityLevel:       newWeightGoal.ActivityLevel,
		DailyCaloriesBudget: newWeightGoal.DailyCalorieBudget,
		CaloriesToMaintain:  newWeightGoal.CaloriesToMaintain,
		Flag:                newWeightGoal.Flag,
		Pace:                newWeightGoal.Pace,
	}, nil
}

func (service *WeightGoalService) GetWeightGoal(ctx context.Context) (*dto.GetWeightGoalResponse, error) {
	profile, err := GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrGetProfile, err)
	}

	wg, err := service.tables.weightGoal.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("profile_id").Eq(profile.ID),
		},
		Sort:  []exp.OrderedExpression{goqu.I("updated_at").Desc()},
		Page:  1,
		Limit: 1,
	})

	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(wg) < 1 {
		return nil, fmt.Errorf("%w", ErrNotFound)
	}

	res := dto.GetWeightGoalResponse{
		StartingWeight:      wg[0].StartingWeight,
		StartingDate:        wg[0].StartingDate.Format(shortdDateLayout),
		TargetWeight:        wg[0].TargetWeight,
		TargetDate:          wg[0].TargetDate.Format(shortdDateLayout),
		ActivityLevel:       wg[0].ActivityLevel,
		DailyCaloriesBudget: wg[0].DailyCalorieBudget,
		CaloriesToMaintain:  wg[0].CaloriesToMaintain,
		Flag:                wg[0].Flag,
		Pace:                wg[0].Pace,
	}

	return &res, nil
}

func (service *WeightGoalService) UpdateWeightGoal(ctx context.Context, body *dto.UpdateWeightGoalRequest) (*dto.CreateWeightGoalResponse, error) {
	var (
		updateWeightGoal = model.WeightGoal{}
		now              = time.Now()
		dateNow          = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startTime, _     = time.Parse(shortdDateLayout, body.StartingDate)
		startDate        = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	)

	// get profile
	profile, err := GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrGetProfile, err)
	}

	// get wg
	wg, err := service.tables.weightGoal.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("profile_id").Eq(profile.ID),
		},
		Sort:  []exp.OrderedExpression{goqu.I("updated_at").Desc()},
		Page:  1,
		Limit: 1,
	})

	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(wg) < 1 {
		return nil, fmt.Errorf("%w", ErrNotFound)
	}

	b, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &updateWeightGoal)
	if err != nil {
		return nil, err
	}

	if startTime.IsZero() {
		startDate = wg[0].StartingDate
	}

	if startDate == dateNow {
		updateWeightGoal.StartingWeight = profile.Weight
	} else if startDate.Before(dateNow) {
		updateWeightGoal.StartingWeight = body.StartingWeight
	}

	if body.TargetWeight > body.TargetWeight {
		updateWeightGoal.Flag = constants.WeightGoalGain
	} else {
		updateWeightGoal.Flag = constants.WeightGoalLoss
	}

	age, _ := strconv.Atoi(profile.Age)
	activLevel := wg[0].ActivityLevel
	if body.ActivityLevel != "" {
		activLevel = body.ActivityLevel
	}

	startWeight := wg[0].StartingWeight
	if body.StartingWeight != 0 {
		startWeight = body.StartingWeight
	}

	targetWeight := wg[0].TargetWeight
	if body.TargetWeight != 0 {
		targetWeight = body.TargetWeight
	}

	pace := wg[0].Pace
	if body.Pace != "" {
		pace = body.Pace
	}

	var wgFlag string
	if body.TargetWeight > body.StartingWeight {
		wgFlag = constants.WeightGoalGain
	} else {
		wgFlag = constants.WeightGoalLoss
	}

	var paceVal float64
	switch strings.ToLower(pace) {
	case constants.WeeklyWeightPaceRelaxed:
		paceVal = constants.WeeklyWeightPaceRelaxedVal
	case constants.WeeklyWeightPaceNormal:
		paceVal = constants.WeeklyWeightPaceNormalVal
	case constants.WeeklyWeightStrict:
		paceVal = constants.WeeklyWeightPaceStrictVal
	}

	caloriesToMaintain := math.Round(CalculateCalorieToMaintain(profile.Sex, startWeight, profile.Height, age, activLevel)*100) / 100
	dailyCaloriesBudget := math.Round(CalculateDailyCalorieBudget(caloriesToMaintain, paceVal, wgFlag)*100) / 100
	targetDate := CalculateTargetDate(profile.Sex, age, profile.Height, startWeight, targetWeight, paceVal, activLevel, wgFlag, startDate)

	updateWeightGoal.CaloriesToMaintain = caloriesToMaintain
	updateWeightGoal.DailyCalorieBudget = dailyCaloriesBudget
	updateWeightGoal.TargetDate = targetDate

	// update wg
	err = service.tables.weightGoal.Update(ctx, wg[0].ID, &updateWeightGoal)
	if err != nil {
		return nil, err
	}

	updatedWG, err := service.tables.weightGoal.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("profile_id").Eq(profile.ID),
		},
		Sort:  []exp.OrderedExpression{goqu.I("updated_at").Desc()},
		Page:  1,
		Limit: 1,
	})

	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	// update profile
	if err := UpdateProfile(ctx, profile.UserID, dto.RequestUpdateProfile{
		Weight:        body.CurrentWeight,
		ActivityLevel: updatedWG[0].ActivityLevel,
	}); err != nil {
		return nil, fmt.Errorf("%w; %w", ErrUpdateProfile, err)
	}

	res := dto.CreateWeightGoalResponse{
		StartingWeight:      updatedWG[0].StartingWeight,
		StartingDate:        updatedWG[0].StartingDate.Format(shortdDateLayout),
		TargetWeight:        updatedWG[0].TargetWeight,
		TargetDate:          updatedWG[0].TargetDate.Format(shortdDateLayout),
		ActivityLevel:       updatedWG[0].ActivityLevel,
		DailyCaloriesBudget: dailyCaloriesBudget,
		CaloriesToMaintain:  caloriesToMaintain,
		Flag:                updatedWG[0].Flag,
		Pace:                pace,
	}

	return &res, nil
}

func (service *WeightGoalService) WightGoalSimulation(ctx context.Context, body dto.SimulationWeightGoalRequest) (*dto.SimulationWeightGoalResponse, error) {
	var (
		now = time.Now()
		res = dto.SimulationWeightGoalResponse{}
	)

	profile, err := GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrGetProfile, err)
	}

	var wgFlag string
	if body.TargetWeight > body.StartingWeight {
		wgFlag = constants.WeightGoalGain
	} else {
		wgFlag = constants.WeightGoalLoss
	}

	age, _ := strconv.Atoi(profile.Age)
	caloriesToMaintain := math.Round(CalculateCalorieToMaintain(profile.Sex, body.StartingWeight, profile.Height, age, body.ActivityLevel)*100) / 100

	paces := map[string]float64{
		constants.WeeklyWeightPaceRelaxed: constants.WeeklyWeightPaceRelaxedVal,
		constants.WeeklyWeightPaceNormal:  constants.WeeklyWeightPaceNormalVal,
		constants.WeeklyWeightStrict:      constants.WeeklyWeightPaceStrictVal,
	}

	for name, val := range paces {
		dailyCaloriesBudget := math.Round(CalculateDailyCalorieBudget(caloriesToMaintain, val, wgFlag)*100) / 100
		targetDate := CalculateTargetDate(profile.Sex, age, profile.Height, body.StartingWeight, body.TargetWeight, val, body.ActivityLevel, wgFlag, now)

		pace := dto.WeightGoalPace{
			Pace:                name,
			DailyCaloriesBudget: dailyCaloriesBudget,
			TargetDate:          targetDate.Format(shortdDateLayout),
		}
		res.Pacing = append(res.Pacing, pace)
	}

	res.StartingWeight = body.StartingWeight
	res.StartingDate = now.Format(shortdDateLayout)
	res.TargetWeight = body.TargetWeight
	res.ActivityLevel = body.ActivityLevel
	res.CaloriesToMaintain = caloriesToMaintain
	res.Flag = wgFlag

	return &res, nil
}
