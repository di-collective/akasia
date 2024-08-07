package api

import (
	"encoding/json"
	"io"
	"monorepo/internal/config"
	"monorepo/internal/constants"
	"monorepo/internal/dto"
	"monorepo/services/calendar/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"
	"github.com/go-playground/validator/v10"

	"github.com/gorilla/schema"
)

type REST struct {
	Router          *chi.Mux
	decoder         *schema.Decoder
	eventService    *service.EventService
	env             *config.Environment
	oauthAuthorizer func(next http.Handler) http.Handler
}

func NewREST(
	eventService *service.EventService,
	env *config.Environment,
) *REST {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(90 * time.Second))
	r.Use(middleware.Compress(6))

	return &REST{
		Router:          r,
		decoder:         schema.NewDecoder(),
		eventService:    eventService,
		env:             env,
		oauthAuthorizer: oauth.Authorize(env.JWTSecret, nil),
	}
}

func (rest *REST) InitializeRoutes() {
	rest.Router.Get("/", rest.Healthcheck)
	rest.Router.Group(func(r chi.Router) {
		r.Use(rest.oauthAuthorizer)
		r.Post("/event", rest.CreateEvent)
		r.Get("/events", rest.GetEvents)
		r.Get("/appointments", rest.GetAppointments)
	})
}

func (rest *REST) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.Object[any]{Data: nil, Message: "OK"})
}

func (rest *REST) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	var req dto.RequestCreateEvent
	json.Unmarshal(payload, &req)

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		err = err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	err = req.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	var profile *dto.ResponseGetProfile
	if req.Type == constants.Appointment {
		code, prof, err := rest.eventService.GetProfile(ctx)
		if err != nil {
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Create Event"})
			return
		}
		profile = prof
	}

	data, err := rest.eventService.CreateEvent(ctx, &req, profile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Create Event"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[[]dto.ResponseCreateEvent]{Data: &data, Message: "OK"})
}

func (rest *REST) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)

	locationID := r.URL.Query().Get("location_id")
	if locationID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "Location id can't be empty", Message: "Invalid query parameter"})
		return
	}

	startTimeStr := r.URL.Query().Get("start_time")
	if startTimeStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "Start time can't be empty", Message: "Invalid query parameter"})
		return
	}
	startTime, err := time.Parse(time.RFC3339, strings.ReplaceAll(startTimeStr, " ", "+"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Invalid query parameter"})
		return
	}

	endTimeStr := r.URL.Query().Get("end_time")
	if endTimeStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "End time can't be empty", Message: "Invalid query parameter"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, strings.ReplaceAll(endTimeStr, " ", "+"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Invalid query parameter"})
		return
	}

	if endTime.Before(startTime) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "End time must be greater than or equal to start time", Message: "Invalid query parameter"})
		return
	}

	code, location, err := rest.eventService.GetLocation(ctx, locationID)
	if err != nil {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Events"})
		return
	}

	code, clinic, err := rest.eventService.GetClinic(ctx, location.ClinicID)
	if err != nil {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Events"})
		return
	}

	data, err := rest.eventService.GetEvents(ctx, dto.FilterGetEvents{
		Page:       page,
		Limit:      limit,
		LocationID: locationID,
		StartTime:  startTime,
		EndTime:    endTime,
	}, location, clinic)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Events"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[dto.ResponseGetEvents]{Data: &data, Message: "OK"})
}

func (rest *REST) GetAppointments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)

	code, profile, err := rest.eventService.GetProfile(ctx)
	if err != nil {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Create Event"})
		return
	}

	data, err := rest.eventService.GetAppointments(ctx, dto.FilterGetAppointments{
		Page:  page,
		Limit: limit,
	}, profile)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Events"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[[]dto.ResponseDetailEvent]{Data: &data, Message: "OK"})
}
