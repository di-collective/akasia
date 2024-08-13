package service

import (
	"context"
	"database/sql"
	"fmt"
	"monorepo/internal/constants"
	"monorepo/internal/dto"
	"monorepo/internal/repository"
	"monorepo/pkg/common"
	"monorepo/services/fitness/model"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/oklog/ulid/v2"
)

func (service *WeightGoalService) IsWeightHistoryExists(ctx context.Context, profileID string, date time.Time) (bool, *model.WeightHistory, error) {
	data := model.WeightHistory{}
	existing, err := service.tables.weightHistory.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("profile_id").Eq(profileID),
			goqu.L("DATE(created_at)").Eq(date.Format(shortdDateLayout)),
		},
		Page:  1,
		Limit: 1,
	})

	if err != nil {
		return false, nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(existing) > 0 {
		data = *existing[0]
	}

	return len(existing) > 0, &data, nil
}

func (service *WeightGoalService) PutWeightHistory(ctx context.Context, body dto.CreateWeightHistoryRequest) (*dto.WeightHistoryResponse, error) {
	var (
		now        = time.Now()
		weightDate time.Time
	)

	if body.Date == "" {
		weightDate = now
	} else {
		weightDate, _ = time.Parse(shortdDateLayout, body.Date)
	}

	profile, err := service.profileService.GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrGetProfile, err)
	}

	// check exist
	isExist, whExist, err := service.IsWeightHistoryExists(ctx, profile.ID, weightDate)
	if err != nil {
		return nil, err
	}

	newWeightHistory := &model.WeightHistory{
		ProfileID: profile.ID,
		Weight:    body.Weight,
		CreatedAt: weightDate,
	}

	if isExist {
		// update weight history
		whExist.Weight = body.Weight
		whExist.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
		err = service.tables.weightHistory.Update(ctx, whExist.ID, whExist)
		if err != nil {
			return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
		}

	} else {
		// insert weight history
		newWeightHistory.ID = ulid.Make().String()
		if err := service.tables.weightHistory.Create(ctx, newWeightHistory); err != nil {
			return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
		}
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

	// update wg to maintain
	if len(wg) > 0 {
		goal := wg[0]
		if IsAchieveGoal(body.Weight, goal.TargetWeight, goal.Flag) && (weightDate.After(goal.TargetDate) || weightDate.Format(shortdDateLayout) == goal.TargetDate.Format(shortdDateLayout)) {
			updateWeightGoal := model.WeightGoal{
				Flag:      constants.WeightGoalMaintain,
				UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			}

			// update wg
			err = service.tables.weightGoal.Update(ctx, wg[0].ID, &updateWeightGoal)
			if err != nil {
				return nil, err
			}
		}
	}

	return &dto.WeightHistoryResponse{
		Weight: newWeightHistory.Weight,
		Date:   newWeightHistory.CreatedAt.Format(shortdDateLayout),
	}, nil
}

func (service *WeightGoalService) GetWeightHistory(ctx context.Context, filter dto.FilterGetWeightHistory) ([]dto.WeightHistoryResponse, error) {
	var res []dto.WeightHistoryResponse

	profile, err := service.profileService.GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrGetProfile, err)
	}

	whFilter := []exp.Expression{goqu.C("profile_id").Eq(profile.ID)}

	if filter.IsCurrent {
		filter.Page = 1
		filter.Limit = 1
	} else if filter.DateFrom != "" || filter.DateTo != "" {
		whFilter = append(whFilter,
			goqu.L("DATE(created_at)").Gte(filter.DateFrom),
			goqu.L("DATE(created_at)").Lte(filter.DateTo),
		)
	}

	wHistories, err := service.tables.weightHistory.List(ctx, &common.FilterOptions{
		Filter: whFilter,
		Sort:   []exp.OrderedExpression{goqu.I("created_at").Desc()},
		Page:   filter.Page,
		Limit:  filter.Limit,
	})

	if err != nil {
		return []dto.WeightHistoryResponse{}, fmt.Errorf("%w; %w", repository.ErrRepositoryQueryFail, err)
	}

	for _, wh := range wHistories {
		data := dto.WeightHistoryResponse{
			Weight: wh.Weight,
			Date:   wh.CreatedAt.Format(shortdDateLayout),
		}

		res = append(res, data)
	}

	return res, nil
}
