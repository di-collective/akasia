package api

import (
	"encoding/json"
	"io"
	"monorepo/internal/config"
	"monorepo/internal/constants"
	"monorepo/internal/dto"
	"monorepo/services/event/service"
	"net/http"
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
