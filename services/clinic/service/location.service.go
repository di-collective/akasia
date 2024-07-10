package service

import (
	"context"
	"fmt"
	"monorepo/internal/dto"
	"monorepo/pkg/common"
	"monorepo/services/clinic/models"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/oklog/ulid/v2"
)

const (
	dateLayout = "2006-01-02T15:04:05Z"
	timeLayout = "15:04:05"
)

func (service *CLinicService) IsLocationExists(ctx context.Context, clinicID, name string) (bool, error) {
	existing, err := service.tables.location.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("clinic_id").Eq(clinicID),
			goqu.C("name").Eq(name),
			goqu.C("deleted_at").IsNull(),
		},
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return false, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	return len(existing) > 0, nil
}

func (service *CLinicService) CreateLocation(ctx context.Context, body *dto.RequestCreateLocation) (*dto.ResponseCreateLocation, error) {
	newLocation := &models.Location{
		ID:          ulid.Make().String(),
		ClinicID:    body.ClinicID,
		Name:        body.Name,
		Address:     body.Address,
		Phone:       body.Phone,
		Capacity:    body.Capacity,
		OpeningTime: body.OpeningTime,
		ClosingTime: body.ClosingTime,
		CreatedAt:   time.Now(),
	}

	isClinicExist, err := service.IsClinicExistsByID(ctx, newLocation.ClinicID)
	if err != nil {
		return nil, err
	}

	if !isClinicExist {
		return nil, fmt.Errorf("%w", ErrClinicNotFound)
	}

	isLocExist, err := service.IsLocationExists(ctx, newLocation.ClinicID, newLocation.Name)
	if err != nil {
		return nil, err
	}

	if isLocExist {
		return nil, fmt.Errorf("%w", ErrLocationExist)
	}

	err = service.tables.location.Create(ctx, newLocation)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	res := dto.ResponseCreateLocation{
		ID:          newLocation.ID,
		ClinicID:    newLocation.ClinicID,
		Name:        newLocation.Name,
		Address:     newLocation.Address,
		Phone:       newLocation.Phone,
		OpeningTime: newLocation.OpeningTime,
		ClosingTime: newLocation.ClosingTime,
	}

	return &res, nil
}

func (service *CLinicService) UpdateLocation(ctx context.Context, locID string, body *dto.RequestUpdateLocation) (*dto.ResponseUpdateLocation, error) {
	location, err := service.tables.location.Get(ctx, locID)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	isClinicExist, err := service.IsClinicExistsByID(ctx, body.ClinicID)
	if err != nil {
		return nil, err
	}

	if !isClinicExist {
		return nil, fmt.Errorf("%w", ErrClinicNotFound)
	}

	isLocExist, err := service.IsLocationExists(ctx, body.ClinicID, body.Name)
	if err != nil {
		return nil, err
	}

	if isLocExist && body.Name != location.Name {
		return nil, fmt.Errorf("%w", ErrLocationExist)
	}

	location.ClinicID = body.ClinicID
	location.Name = body.Name
	location.Address = body.Address
	location.Phone = body.Phone
	location.OpeningTime = body.OpeningTime
	location.ClosingTime = body.ClosingTime

	err = service.tables.location.Update(ctx, locID, location)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	res := dto.ResponseUpdateLocation{
		ID:          locID,
		ClinicID:    location.ClinicID,
		Name:        location.Name,
		Address:     location.Address,
		Phone:       location.Phone,
		OpeningTime: location.OpeningTime,
		ClosingTime: location.ClosingTime,
		Capacity:    location.Capacity,
	}

	return &res, nil
}

func (service *CLinicService) DeleteLocation(ctx context.Context, id string) error {
	location, err := service.tables.location.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	err = service.tables.location.Delete(ctx, location.ID)
	if err != nil {
		return fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	return nil
}

func (service *CLinicService) GetLocation(ctx context.Context, id string) (*dto.ResponseGetLocation, error) {
	location, err := service.tables.location.Get(ctx, id)
	if err != nil {
		if err == ErrNoResult {
			return nil, fmt.Errorf("%w", "data not found")
		}
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	openTime, _ := time.Parse(dateLayout, location.OpeningTime)
	closeTime, _ := time.Parse(dateLayout, location.ClosingTime)

	res := dto.ResponseGetLocation{}
	res.ID = location.ID
	res.ClinicID = location.ClinicID
	res.Name = location.Name
	res.Address = location.Address
	res.Phone = location.Phone
	res.Capacity = location.Capacity
	res.OpeningTime = openTime.Format(timeLayout)
	res.ClosingTime = closeTime.Format(timeLayout)
	res.CreatedAt = location.CreatedAt

	if location.DeletedAt.Valid {
		res.DeletedAt = &location.DeletedAt.Time
	}

	return &res, nil
}

func (service *CLinicService) GetLocationByClinic(ctx context.Context, clinicID string) ([]dto.ResponseGetLocation, error) {
	res := []dto.ResponseGetLocation{}

	locations, err := service.tables.location.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("clinic_id").Eq(clinicID),
		},
		Sort: []exp.OrderedExpression{goqu.I("id").Desc()},
	})

	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	for _, loc := range locations {
		c := dto.ResponseGetLocation{}
		c.ID = loc.ID
		c.ClinicID = loc.ClinicID
		c.Name = loc.Name
		c.Address = loc.Address
		c.Phone = loc.Phone
		c.OpeningTime = loc.OpeningTime
		c.ClosingTime = loc.ClosingTime
		c.CreatedAt = loc.CreatedAt
		if loc.DeletedAt.Valid {
			c.DeletedAt = &loc.DeletedAt.Time
		}

		res = append(res, c)
	}

	return res, nil
}
