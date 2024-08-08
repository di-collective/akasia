package service

import (
	"context"
	"errors"
	"fmt"
	"monorepo/internal/dto"
	"monorepo/pkg/utils"
	"net/http"
	"os"

	"github.com/go-chi/oauth"
)

func GetProfile(ctx context.Context) (*dto.ResponseGetProfile, error) {
	url := os.Getenv("USER_BASE_URL") + "/profile"

	type GetProfileResponse struct {
		Data    dto.ResponseGetProfile `json:"data"`
		Message string                 `json:"message"`
	}
	resp := GetProfileResponse{}

	headers := []utils.Header{
		{
			Key:   "Authorization",
			Value: "Bearer " + ctx.Value(oauth.AccessTokenContext).(string),
		},
	}

	if _, err := utils.DoRequest(http.MethodGet, url, headers, nil, &resp); err != nil {
		errMessage := fmt.Sprintf("call api error : %s", err.Error())
		return nil, errors.New(errMessage)
	}

	return &resp.Data, nil
}

func UpdateProfile(ctx context.Context, userID string, data dto.RequestUpdateProfile) error {
	url := os.Getenv("USER_BASE_URL") + "/profile/" + userID
	var resp any

	headers := []utils.Header{
		{
			Key:   "Authorization",
			Value: "Bearer " + ctx.Value(oauth.AccessTokenContext).(string),
		},
	}

	if _, err := utils.DoRequest(http.MethodPatch, url, headers, data, &resp); err != nil {
		errMessage := fmt.Sprintf("call api error : %s", err.Error())
		return errors.New(errMessage)
	}

	return nil
}
