package service

import (
	"context"
	"fmt"
	"math"
	"monorepo/internal/constants"
	"monorepo/internal/dto"
	"monorepo/pkg/common"
	"monorepo/services/medical-record/model"
	"strconv"
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
	if body.TargetWeight > body.TargetWeight {
		wgFlag = constants.WeightGoalGain
	} else {
		wgFlag = constants.WeightGoalLoss
	}

	// calculate calorie budget
	age, _ := strconv.Atoi(profile.Age)
	var calorieBgt = math.Round(CalculateCalorieBudget(profile.Sex, body.StartingWeight, profile.Height, age, body.ActivityLevel)*100) / 100

	//calculate target date
	targetDate := CalculateTargetDate(profile.Sex, age, profile.Height, body.StartingWeight, body.TargetWeight, 0.5, body.ActivityLevel, wgFlag, now)

	newWeightGoal := &model.WeightGoal{
		ID:             ulid.Make().String(),
		ProfileID:      profile.ID,
		StartingWeight: body.StartingWeight,
		StartingDate:   now,
		TargetWeight:   body.TargetWeight,
		TargetDate:     targetDate,
		CalorieBudget:  calorieBgt,
		Flag:           wgFlag,
		CreatedAt:      now,
		ActivityLevel:  body.ActivityLevel,
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
		StartingWeight: newWeightGoal.StartingWeight,
		StartingDate:   newWeightGoal.StartingDate.Format(shortdDateLayout),
		TargetWeight:   newWeightGoal.TargetWeight,
		TargetDate:     newWeightGoal.TargetDate.Format(shortdDateLayout),
		ActivityLevel:  newWeightGoal.ActivityLevel,
		CalorieBudget:  newWeightGoal.CalorieBudget,
		Flag:           newWeightGoal.Flag,
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
		StartingWeight: wg[0].StartingWeight,
		StartingDate:   wg[0].StartingDate.Format(shortdDateLayout),
		TargetWeight:   wg[0].TargetWeight,
		TargetDate:     wg[0].TargetDate.Format(shortdDateLayout),
		ActivityLevel:  wg[0].ActivityLevel,
		CalorieBudget:  wg[0].CalorieBudget,
		Flag:           wg[0].Flag,
	}

	return &res, nil
}
