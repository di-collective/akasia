package api

import (
	"encoding/json"
	"io"
	"monorepo/internal/config"
	"monorepo/internal/dto"
	"monorepo/services/fitness/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type REST struct {
	Router *chi.Mux

	decoder         *schema.Decoder
	clinicService   *service.WeightGoalService
	env             *config.Environment
	oauthAuthorizer func(next http.Handler) http.Handler
}

func NewREST(
	weightGoalService *service.WeightGoalService,
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
		clinicService:   weightGoalService,
		env:             env,
		oauthAuthorizer: oauth.Authorize(env.JWTSecret, nil),
	}
}

func (rest *REST) InitializeRoutes() {
	rest.Router.Get("/", rest.Healthcheck)
	rest.Router.Group(func(r chi.Router) {
		r.Use(rest.oauthAuthorizer)
		r.Post("/weight-goal", rest.CreateWeightGoal)
		r.Get("/weight-goal", rest.GetWeightGoal)
		r.Patch("/weight-goal", rest.UpdateWeightGoal)
		r.Post("/weight-goal/simulation", rest.WeightGoalSimulation)
	})
}

func (rest *REST) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.Object[any]{Data: nil, Message: "OK"})
}

func (rest *REST) CreateWeightGoal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	claims := ctx.Value(oauth.ClaimsContext)
	c, _ := json.Marshal(claims)
	var fClaims dto.FirebaseClaims
	json.Unmarshal(c, &fClaims)

	var req dto.CreateWeightGoalRequest
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

	data, err := rest.clinicService.CreateWightGoal(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to create weight goal"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.CreateWeightGoalResponse]{Data: &data, Message: "OK"})
}

func (rest *REST) GetWeightGoal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := rest.clinicService.GetWeightGoal(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to get weight goal"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.GetWeightGoalResponse]{Data: &data, Message: "OK"})
}

func (rest *REST) UpdateWeightGoal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	claims := ctx.Value(oauth.ClaimsContext)
	c, _ := json.Marshal(claims)
	var fClaims dto.FirebaseClaims
	json.Unmarshal(c, &fClaims)

	var req dto.UpdateWeightGoalRequest
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

	data, err := rest.clinicService.UpdateWeightGoal(ctx, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to update weight goal"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.CreateWeightGoalResponse]{Data: &data, Message: "OK"})
}

func (rest *REST) WeightGoalSimulation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	claims := ctx.Value(oauth.ClaimsContext)
	c, _ := json.Marshal(claims)
	var fClaims dto.FirebaseClaims
	json.Unmarshal(c, &fClaims)

	var req dto.SimulationWeightGoalRequest
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

	data, err := rest.clinicService.WightGoalSimulation(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to simulate weight goal"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.SimulationWeightGoalResponse]{Data: &data, Message: "OK"})
}
