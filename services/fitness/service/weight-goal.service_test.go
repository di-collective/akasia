package service

import (
	"context"
	"errors"
	"monorepo/internal/dto"
	"reflect"
	"testing"
	"time"

	mock "monorepo/services/fitness/service/mock"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/go-chi/oauth"
	"github.com/golang/mock/gomock"
)

func TestWeightGoalService_WightGoalSimulation(t *testing.T) {
	mockTime := time.Date(2024, time.August, 12, 0, 0, 0, 0, time.UTC)
	patches := gomonkey.ApplyFunc(time.Now, func() time.Time {
		return mockTime
	})
	defer patches.Reset()

	ctx := context.WithValue(context.Background(), oauth.AccessTokenContext, "")

	type args struct {
		ctx  context.Context
		body dto.SimulationWeightGoalRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMocks func(ctrl *gomock.Controller) ProfileServiceInterface
		want       *dto.SimulationWeightGoalResponse
		wantErr    bool
	}{
		{
			name: "Response is match with sadentary level",
			args: args{
				ctx: ctx,
				body: dto.SimulationWeightGoalRequest{
					StartingWeight: 80,
					TargetWeight:   60,
					ActivityLevel:  "Sadentary",
				},
			},
			setupMocks: func(ctrl *gomock.Controller) ProfileServiceInterface {
				mockProfile := mock.NewMockProfileServiceInterface(ctrl)
				mockProfile.EXPECT().GetProfile(ctx).Return(&dto.ResponseGetProfile{
					Age:    "30",
					Height: 160,
					Sex:    "Male",
				}, nil)
				return mockProfile
			},
			want: &dto.SimulationWeightGoalResponse{
				StartingWeight:     80,
				StartingDate:       time.Now().Format(shortdDateLayout),
				TargetWeight:       60,
				ActivityLevel:      "Sadentary",
				CaloriesToMaintain: 2117.38,
				Flag:               "loss",
				Pacing: []dto.WeightGoalPace{
					{
						Pace:                "relaxed",
						DailyCaloriesBudget: 1842.38,
						TargetDate:          "2026-02-23",
					},
					{
						Pace:                "normal",
						DailyCaloriesBudget: 1567.38,
						TargetDate:          "2025-05-19",
					},
					{
						Pace:                "strict",
						DailyCaloriesBudget: 1017.38,
						TargetDate:          "2024-12-30",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Response with gain flag",
			args: args{
				ctx: context.WithValue(context.Background(), oauth.AccessTokenContext, ""),
				body: dto.SimulationWeightGoalRequest{
					StartingWeight: 60,
					TargetWeight:   70,
					ActivityLevel:  "Sadentary",
				},
			},
			setupMocks: func(ctrl *gomock.Controller) ProfileServiceInterface {
				mockProfile := mock.NewMockProfileServiceInterface(ctrl)
				mockProfile.EXPECT().GetProfile(ctx).Return(&dto.ResponseGetProfile{
					Age:    "30",
					Height: 160,
					Sex:    "Male",
				}, nil)
				return mockProfile
			},
			want: &dto.SimulationWeightGoalResponse{
				StartingWeight:     60,
				StartingDate:       time.Now().Format(shortdDateLayout),
				TargetWeight:       70,
				ActivityLevel:      "Sadentary",
				CaloriesToMaintain: 1787.38,
				Flag:               "gain",
				Pacing: []dto.WeightGoalPace{
					{
						Pace:                "relaxed",
						DailyCaloriesBudget: 2062.38,
						TargetDate:          "2025-05-19",
					},
					{
						Pace:                "normal",
						DailyCaloriesBudget: 2337.38,
						TargetDate:          "2024-12-30",
					},
					{
						Pace:                "strict",
						DailyCaloriesBudget: 2887.38,
						TargetDate:          "2024-10-21",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error get profile",
			args: args{
				ctx: context.WithValue(context.Background(), oauth.AccessTokenContext, ""),
				body: dto.SimulationWeightGoalRequest{
					StartingWeight: 80,
					TargetWeight:   60,
					ActivityLevel:  "Sadentary",
				},
			},
			setupMocks: func(ctrl *gomock.Controller) ProfileServiceInterface {
				mockProfile := mock.NewMockProfileServiceInterface(ctrl)
				mockProfile.EXPECT().GetProfile(ctx).Return(nil, errors.New("error get profile"))
				return mockProfile
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProfile := tt.setupMocks(ctrl)

			service := &WeightGoalService{
				profileService: mockProfile,
			}
			got, err := service.WightGoalSimulation(tt.args.ctx, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("WeightGoalService.WightGoalSimulation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WeightGoalService.WightGoalSimulation() = %v, want %v", got, tt.want)
			}
		})
	}
}
