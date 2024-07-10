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
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

func NewClinicService(
	tbCLinic common.Repository[models.Clinic, string],
	tbLocation common.Repository[models.Location, string],
) *CLinicService {
	service := &CLinicService{}
	service.validate = validator.New()
	service.tables.clinic = tbCLinic
	service.tables.location = tbLocation

	return service
}

type CLinicService struct {
	validate *validator.Validate
	tables   struct {
		clinic   common.Repository[models.Clinic, string]
		location common.Repository[models.Location, string]
	}
}

func (service *CLinicService) IsClinicExists(ctx context.Context, name string) (bool, error) {
	existing, err := service.tables.clinic.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
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

func (service *CLinicService) IsClinicExistsByID(ctx context.Context, id string) (bool, error) {
	existing, err := service.tables.clinic.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("id").Eq(id),
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

func (service *CLinicService) CreateClinic(ctx context.Context, body *dto.RequestCreateClinic) (*dto.ResponseCreateClinic, error) {
	newClinic := &models.Clinic{
		ID:        ulid.Make().String(),
		Name:      body.Name,
		Address:   body.Address,
		Phone:     body.Phone,
		Logo:      body.Logo,
		CreatedAt: time.Now(),
	}

	isExist, err := service.IsClinicExists(ctx, newClinic.Name)
	if err != nil {
		return nil, err
	}

	if isExist {
		return nil, fmt.Errorf("%w", ErrClinicExist)
	}

	err = service.tables.clinic.Create(ctx, newClinic)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	res := dto.ResponseCreateClinic{
		ID:      newClinic.ID,
		Name:    newClinic.Name,
		Address: newClinic.Address,
		Phone:   newClinic.Phone,
		Logo:    newClinic.Logo,
	}

	return &res, nil
}

func (service *CLinicService) UpdateClinic(ctx context.Context, id string, body *dto.RequestUpdateClinic) (*dto.ResponseUpdateClinic, error) {
	clinic, err := service.tables.clinic.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	isExist, err := service.IsClinicExists(ctx, body.Name)
	if err != nil {
		return nil, err
	}

	if isExist && body.Name != clinic.Name {
		return nil, fmt.Errorf("%w", ErrClinicExist)
	}

	clinic.Name = body.Name
	clinic.Address = body.Address
	clinic.Phone = body.Phone
	clinic.Logo = body.Logo

	err = service.tables.clinic.Update(ctx, clinic.ID, clinic)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	res := dto.ResponseUpdateClinic{
		ID:      clinic.ID,
		Name:    clinic.Name,
		Address: clinic.Address,
		Phone:   clinic.Phone,
		Logo:    clinic.Logo,
	}

	return &res, nil
}

func (service *CLinicService) DeleteClinic(ctx context.Context, id string) error {
	clinic, err := service.tables.clinic.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	err = service.tables.clinic.Delete(ctx, clinic.ID)
	if err != nil {
		return fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	return nil
}

func (service *CLinicService) GetClinic(ctx context.Context, id string) (*dto.ResponseGetClinic, error) {
	clinic, err := service.tables.clinic.Get(ctx, id)
	if err != nil {
		if err == ErrNoResult {
			return nil, fmt.Errorf("%w", "data not found")
		}
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	res := dto.ResponseGetClinic{}
	res.ID = clinic.ID
	res.Name = clinic.Name
	res.Address = clinic.Address
	res.Phone = clinic.Phone
	res.Logo = clinic.Logo
	res.CreatedAt = clinic.CreatedAt

	if clinic.DeletedAt.Valid {
		res.DeletedAt = &clinic.DeletedAt.Time
	}

	return &res, nil
}

func (service *CLinicService) GetAllClinic(ctx context.Context, filter dto.FilterGetClinic) ([]dto.ResponseGetClinic, error) {
	res := []dto.ResponseGetClinic{}

	clinics, err := service.tables.clinic.List(ctx, &common.FilterOptions{
		Sort:  []exp.OrderedExpression{goqu.I("id").Desc()},
		Page:  filter.Page,
		Limit: filter.Limit,
	})

	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	for _, clinic := range clinics {
		c := dto.ResponseGetClinic{}
		c.ID = clinic.ID
		c.Name = clinic.Name
		c.Address = clinic.Address
		c.Phone = clinic.Phone
		c.Logo = clinic.Logo
		c.CreatedAt = clinic.CreatedAt

		if clinic.DeletedAt.Valid {
			c.DeletedAt = &clinic.DeletedAt.Time
		}

		res = append(res, c)
	}

	return res, nil
}
