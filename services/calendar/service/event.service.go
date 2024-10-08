package service

import (
	"context"
	"encoding/json"
	"fmt"
	"monorepo/internal/dto"
	"monorepo/pkg/common"
	"monorepo/pkg/utils"
	"monorepo/services/calendar/models"
	"os"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/go-chi/oauth"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"

	"monorepo/internal/constants"
	"monorepo/internal/repository"
)

func NewEventService(tbEvent common.Repository[models.Event, string]) *EventService {
	service := &EventService{}
	service.validate = validator.New()
	service.tables.event = tbEvent
	return service
}

type EventService struct {
	validate *validator.Validate
	tables   struct {
		event common.Repository[models.Event, string]
	}
}

func (service *EventService) GetProfile(ctx context.Context) (int, *dto.ResponseGetProfile, error) {
	method := "GET"

	url := os.Getenv("BASE_URL_USER") + "/profile"

	headers := []utils.Header{
		{
			Key:   "Authorization",
			Value: "Bearer " + ctx.Value(oauth.AccessTokenContext).(string),
		},
	}

	profile := dto.ResponseGetProfile{}

	var result map[string]interface{}

	res, err := utils.DoRequest(method, url, headers, nil, &result)
	if err != nil {
		return res.StatusCode, &profile, err
	}

	err = json.Unmarshal([]byte(res.Body), &result)
	if err != nil {
		return res.StatusCode, &profile, err
	}

	data := result["data"].(map[string]interface{})

	profileData, err := json.Marshal(data)
	if err != nil {
		return res.StatusCode, &profile, err
	}

	err = json.Unmarshal(profileData, &profile)
	if err != nil {
		return res.StatusCode, &profile, err
	}

	return res.StatusCode, &profile, nil
}

func (service *EventService) IsEventExists(ctx context.Context, locationID string, startTime time.Time, endTime time.Time, _type string) (bool, error) {
	var filters []exp.Expression
	var existing []*models.Event
	var err error
	if _type == constants.Holiday {
		filters = append(filters,
			goqu.C("location_id").Eq(locationID),
			goqu.And(
				goqu.C("start_time").Eq(startTime),
				goqu.C("end_time").Eq(endTime),
			),
			goqu.C("deleted_at").IsNull(),
		)
		existing, err = service.tables.event.List(ctx, &common.FilterOptions{
			Filter: filters,
			Page:   1,
			Limit:  1,
		})
		if err != nil {
			return false, fmt.Errorf("%w; %w", repository.ErrRepositoryQueryFail, err)
		}
	}

	return len(existing) > 0, nil
}

func (service *EventService) CreateEvent(ctx context.Context, body *dto.RequestCreateEvent, profile *dto.ResponseGetProfile) ([]dto.ResponseCreateEvent, error) {
	var res []dto.ResponseCreateEvent
	var events []models.Event

	if body.Type == constants.Holiday {
		events = service.createHolidayEvents(body)
	} else {
		event, err := service.createSingleEvent(body, profile)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	for _, event := range events {
		isExist, err := service.IsEventExists(ctx, event.LocationID, event.StartTime, event.EndTime, event.Type)
		if err != nil {
			if event.Type == constants.Holiday {
				errMsg := fmt.Errorf("%w; %w", repository.ErrRepositoryMutateFail, err)
				logrus.Errorf("%v", errMsg)
				continue
			}
			return nil, fmt.Errorf("%w; %w", repository.ErrRepositoryMutateFail, err)
		}

		if isExist {
			if event.Type == constants.Holiday {
				warnMsg := fmt.Errorf("%w", repository.ErrExist)
				logrus.Errorf("%v", warnMsg)
				continue
			}
			return nil, fmt.Errorf("%w", repository.ErrExist)
		}

		err = service.tables.event.Create(ctx, &event)
		if err != nil {
			if event.Type == constants.Holiday {
				errMsg := fmt.Errorf("%w; %w", repository.ErrRepositoryMutateFail, err)
				logrus.Errorf("%v", errMsg)
				continue
			}
			return nil, fmt.Errorf("%w; %w", repository.ErrRepositoryMutateFail, err)
		}

		res = append(res, dto.ResponseCreateEvent{
			ID:         event.ID,
			ProfileID:  event.ProfileID,
			LocationID: event.LocationID,
			Status:     event.Status,
			Type:       event.Type,
			StartTime:  event.StartTime,
			EndTime:    event.EndTime,
		})
	}

	return res, nil
}

func (service *EventService) createSingleEvent(body *dto.RequestCreateEvent, profile *dto.ResponseGetProfile) (models.Event, error) {
	event := models.Event{
		ID:         ulid.Make().String(),
		ProfileID:  &profile.ID,
		LocationID: body.LocationID,
		Type:       body.Type,
		CreatedAt:  time.Now(),
		Status:     body.Status,
		StartTime:  body.StartTime,
		EndTime:    body.StartTime.Add(1 * time.Hour),
	}
	return event, nil
}

func (service *EventService) createHolidayEvents(body *dto.RequestCreateEvent) []models.Event {
	var events []models.Event
	startDate := time.Date(body.StartTime.Year(), body.StartTime.Month(), body.StartTime.Day(), 0, 0, 0, 0, body.StartTime.Location())
	endDate := body.EndTime

	for currentTime := startDate; !currentTime.After(endDate); currentTime = currentTime.Add(24 * time.Hour) {
		events = append(events, models.Event{
			ID:         ulid.Make().String(),
			LocationID: body.LocationID,
			Type:       body.Type,
			CreatedAt:  time.Now(),
			Status:     constants.Scheduled,
			StartTime:  currentTime,
			EndTime:    time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()),
		})
	}

	return events
}

func (service *EventService) GetEvents(ctx context.Context, filter dto.FilterGetEvents, location *dto.ResponseGetLocation, clinic *dto.ResponseGetClinic) (dto.ResponseGetEvents, error) {
	var res dto.ResponseGetEvents

	events, err := service.tables.event.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("location_id").Eq(filter.LocationID),
			goqu.C("start_time").Gte(filter.StartTime),
			goqu.C("end_time").Lte(filter.EndTime),
			goqu.C("deleted_at").IsNull(),
			goqu.C("status").Neq(constants.Canceled),
			goqu.C("type").Eq(filter.Type),
		},
		Sort:  []exp.OrderedExpression{goqu.I("start_time").Asc()},
		Page:  filter.Page,
		Limit: filter.Limit,
	})

	if err != nil {
		return dto.ResponseGetEvents{}, fmt.Errorf("%w; %w", repository.ErrRepositoryQueryFail, err)
	}

	for _, event := range events {
		e := dto.ResponseDetailEvent{
			Status:    event.Status,
			Type:      event.Type,
			StartTime: event.StartTime,
			EndTime:   event.EndTime,
			Clinic:    clinic.Name,
			Location:  location.Name,
		}

		res.Events = append(res.Events, e)
	}

	res.Capacity, _ = strconv.Atoi(os.Getenv("CAPACITY"))

	return res, nil
}

func (service *EventService) GetLocation(ctx context.Context, id string) (int, *dto.ResponseGetLocation, error) {
	method := "GET"

	url := os.Getenv("BASE_URL_CLINIC") + "/clinic/location/" + id

	headers := []utils.Header{
		{
			Key:   "Authorization",
			Value: "Bearer " + ctx.Value(oauth.AccessTokenContext).(string),
		},
	}

	location := dto.ResponseGetLocation{}

	var result map[string]interface{}

	res, err := utils.DoRequest(method, url, headers, nil, &result)
	if err != nil {
		return res.StatusCode, &location, err
	}

	err = json.Unmarshal([]byte(res.Body), &result)
	if err != nil {
		return res.StatusCode, &location, err
	}

	data := result["data"].(map[string]interface{})

	locationData, err := json.Marshal(data)
	if err != nil {
		return res.StatusCode, &location, err
	}

	err = json.Unmarshal(locationData, &location)
	if err != nil {
		return res.StatusCode, &location, err
	}

	return res.StatusCode, &location, nil
}

func (service *EventService) GetClinic(ctx context.Context, id string) (int, *dto.ResponseGetClinic, error) {
	method := "GET"

	url := os.Getenv("BASE_URL_CLINIC") + "/clinic/" + id

	headers := []utils.Header{
		{
			Key:   "Authorization",
			Value: "Bearer " + ctx.Value(oauth.AccessTokenContext).(string),
		},
	}

	clinic := dto.ResponseGetClinic{}

	var result map[string]interface{}

	res, err := utils.DoRequest(method, url, headers, nil, &result)
	if err != nil {
		return res.StatusCode, &clinic, err
	}

	err = json.Unmarshal([]byte(res.Body), &result)
	if err != nil {
		return res.StatusCode, &clinic, err
	}

	data := result["data"].(map[string]interface{})

	clinicData, err := json.Marshal(data)
	if err != nil {
		return res.StatusCode, &clinic, err
	}

	err = json.Unmarshal(clinicData, &clinic)
	if err != nil {
		return res.StatusCode, &clinic, err
	}

	return res.StatusCode, &clinic, nil
}

func (service *EventService) GetAppointments(ctx context.Context, filter dto.FilterGetAppointments, profile *dto.ResponseGetProfile) ([]dto.ResponseDetailEvent, error) {
	var res []dto.ResponseDetailEvent

	events, err := service.tables.event.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("profile_id").Eq(profile.ID),
			goqu.C("type").Eq(constants.Appointment),
			goqu.C("deleted_at").IsNull(),
		},
		Sort:  []exp.OrderedExpression{goqu.I("start_time").Desc()},
		Page:  filter.Page,
		Limit: filter.Limit,
	})

	if err != nil {
		return []dto.ResponseDetailEvent{}, fmt.Errorf("%w; %w", repository.ErrRepositoryQueryFail, err)
	}

	for _, event := range events {
		e := dto.ResponseDetailEvent{
			Status:    event.Status,
			Type:      event.Type,
			StartTime: event.StartTime,
			EndTime:   event.EndTime,
		}
		_, location, err := service.GetLocation(ctx, event.LocationID)
		if err == nil && location != nil && location.Name != "" {
			e.Location = location.Name
		}
		_, clinic, err := service.GetClinic(ctx, location.ClinicID)
		if err == nil && clinic != nil && clinic.Name != "" {
			e.Clinic = clinic.Name
		}

		res = append(res, e)
	}

	return res, nil
}
