package service

import (
	"context"
	"fmt"
	"monorepo/internal/dto"
	"monorepo/pkg/common"
	"monorepo/pkg/utils"
	"monorepo/services/user/models"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(
	tbUser common.Repository[models.User, string],
	tbProfile common.Repository[models.Profile, string],
) *UserService {
	service := &UserService{}
	service.validate = validator.New()
	service.tables.user = tbUser
	service.tables.profile = tbProfile

	return service
}

type UserService struct {
	validate *validator.Validate
	tables   struct {
		user    common.Repository[models.User, string]
		profile common.Repository[models.Profile, string]
	}
}

func (service *UserService) IsHandleExists(ctx context.Context, handle string) (bool, error) {
	existing, err := service.tables.user.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{goqu.C("handle").Eq(handle)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return false, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	return len(existing) > 0, nil
}

func (service *UserService) RegisterUser(ctx context.Context, body *dto.RequestRegisterUser) (*models.User, error) {
	existing, err := service.tables.user.List(ctx, &common.FilterOptions{
		Sort:   []exp.Expression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("handle").Eq(body.Email)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(existing) > 0 {
		return existing[0], nil
	}

	if body.Password == "" {
		body.Password = utils.RandAlphanumericString(12)
	}

	plainTextPassword := utils.Ternary(body.Provider == "email", body.Password, ulid.Make().String())
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrPasswordHashingFailed, err)
	}

	newUser := &models.User{
		ID:        ulid.Make().String(),
		Provider:  body.Provider,
		Handle:    body.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}
	err = service.tables.user.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	return newUser, nil
}

func (service *UserService) ChangePassword(ctx context.Context) error { return nil }
func (service *UserService) ResetPassword(ctx context.Context) error  { return nil }
func (service *UserService) UpdateUser(ctx context.Context) error     { return nil }
