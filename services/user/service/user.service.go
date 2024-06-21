package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"monorepo/internal/config"
	"monorepo/internal/dto"
	"monorepo/pkg/common"
	"monorepo/pkg/utils"
	"monorepo/services/user/models"
	"net/url"
	"os"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(
	tbUser common.Repository[models.User, string],
	tbProfile common.Repository[models.Profile, string],
	fbaClient *auth.Client,
) *UserService {
	service := &UserService{}
	service.fbaClient = fbaClient
	service.validate = validator.New()
	service.tables.user = tbUser
	service.tables.profile = tbProfile

	return service
}

type UserService struct {
	fbaClient *auth.Client
	validate  *validator.Validate
	tables    struct {
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
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
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

func (service *UserService) CreateProfile(ctx context.Context, body *dto.RequestCreateProfile) (*dto.ResponseCreateProfile, error) {
	user, err := service.tables.user.List(ctx, &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("id").Eq(body.UserID)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(user) == 0 {
		err = errors.New("user does not exist")
		return nil, err
	}

	existing, err := service.tables.profile.List(ctx, &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("user_id").Eq(body.UserID)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(existing) > 0 {
		err = errors.New("profile already exist")
		return nil, err
	}

	name := strings.Fields(body.Name)

	newProfile := &models.Profile{
		ID:          ulid.Make().String(),
		UserID:      body.UserID,
		MedicalID:   ulid.Make().String(),
		FirstName:   strings.Join(name[:len(name)-1], " "),
		LastName:    name[len(name)-1],
		CountryCode: body.CountryCode,
		Phone:       body.Phone,
		NIK:         &body.NIK,
		PhotoUrl:    &body.PhotoUrl,
		CreatedAt:   time.Now(),
	}
	err = service.tables.profile.Create(ctx, newProfile)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	res := dto.ResponseCreateProfile{
		ID:          newProfile.ID,
		UserID:      newProfile.UserID,
		MedicalID:   newProfile.MedicalID,
		Name:        body.Name,
		CountryCode: newProfile.CountryCode,
		Phone:       newProfile.Phone,
		NIK:         *newProfile.NIK,
		PhotoUrl:    *newProfile.PhotoUrl,
	}

	return &res, nil
}

func (service *UserService) GetProfile(ctx context.Context, body *dto.FirebaseClaims) (any, error) {
	profile, err := service.tables.profile.List(ctx, &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("user_id").Eq(body.UserID)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(profile) == 0 {
		err = errors.New("profile not found")
		return nil, err
	}

	res := dto.ResponseGetProfile{}
	p, err := json.Marshal(profile[0])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(p, &res)
	if err != nil {
		return nil, err
	}

	res.Role = body.Role
	res.Name = fmt.Sprintf("%s %s", profile[0].FirstName, profile[0].LastName)

	return &res, nil
}

func (service *UserService) UpdateProfile(ctx context.Context, userId string, body *dto.RequestUpdateProfile) (any, error) {
	updateProfile := models.Profile{}

	// get profile
	profile, err := service.tables.profile.List(ctx, &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("user_id").Eq(userId)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(profile) == 0 {
		err = errors.New("profile not found")
		return nil, err
	}

	// update profile
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &updateProfile)
	if err != nil {
		return nil, err
	}

	if updateProfile.DOB.IsZero() {
		updateProfile.DOB = profile[0].DOB
	}

	if updateProfile.PhotoUrl != nil {
		user, err := service.tables.user.List(ctx, &common.FilterOptions{
			Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
			Filter: []exp.Expression{goqu.C("id").Eq(userId)},
			Page:   1,
			Limit:  1,
		})
		if err != nil {
			return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
		}

		// get user by email
		u, err := service.fbaClient.GetUserByEmail(ctx, user[0].Handle)
		if err != nil {
			return nil, err
		}

		// update firebase photo
		params := (&auth.UserToUpdate{}).
			PhotoURL(*updateProfile.PhotoUrl)
		_, err = service.fbaClient.UpdateUser(ctx, u.UID, params)
		if err != nil {
			return nil, err
		}
	}

	err = service.tables.profile.Update(ctx, profile[0].ID, &updateProfile)
	if err != nil {
		return nil, err
	}

	// get updated profile
	updatedProfile, err := service.tables.profile.List(ctx, &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("user_id").Eq(userId)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	return updatedProfile[0], nil
}

func (service *UserService) DeleteProfile(ctx context.Context, userId string) error {
	profile, err := service.tables.profile.List(ctx, &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("user_id").Eq(userId)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(profile) == 0 {
		err = errors.New("profile not found")
		return err
	}

	updateProfile := models.Profile{
		DeletedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err = service.tables.profile.Update(ctx, profile[0].ID, &updateProfile)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) UploadPhoto(ctx context.Context, env *config.Environment, file multipart.File, fileName, userId string) (any, error) {
	// create temp file
	tempFile, err := os.CreateTemp("./", fileName)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// read file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFile.Write(fileBytes)

	link := fmt.Sprintf("https://%s.%s/%s", env.OSSBucketName, env.OSSEndpoint, url.QueryEscape(fileName))
	err = utils.PutFile(env.OSSBucketName, env.OSSEndpoint, env.OSSAccessKeyID,
		env.OSSAccessKeySecret, tempFile.Name(), fileName)
	if err != nil {
		return "", err
	}

	// save data to db
	updatePhoto := dto.RequestUpdateProfile{
		PhotoUrl: link,
	}
	_, err = service.UpdateProfile(ctx, userId, &updatePhoto)
	if err != nil {
		return "", err
	}

	// remove temp file
	err = os.Remove(tempFile.Name())
	if err != nil {
		return "", err
	}

	return link, nil
}
