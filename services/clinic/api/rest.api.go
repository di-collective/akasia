package api

import (
	"encoding/json"
	"io"
	"monorepo/internal/config"
	"monorepo/internal/dto"
	"monorepo/services/clinic/service"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"
	"github.com/go-playground/validator/v10"

	"github.com/gorilla/schema"
)

type REST struct {
	Router *chi.Mux

	decoder       *schema.Decoder
	clinicService *service.CLinicService
	env           *config.Environment
}

func NewREST(
	clinicService *service.CLinicService,
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
		Router:        r,
		decoder:       schema.NewDecoder(),
		clinicService: clinicService,
		env:           env,
	}
}

func (rest *REST) InitializeRoutes() {
	rest.Router.Get("/", rest.Healthcheck)
	rest.Router.Group(func(r chi.Router) {
		// r.Use(rest.oauthAuthorizer)

		r.Post("/clinic", rest.CreateClinic)
		r.Get("/clinic", rest.GetAllClinic)
		r.Get("/clinic/{id}", rest.GetClinic)
		r.Patch("/clinic/{id}", rest.UpdateClinic)
		r.Delete("/clinic/{id}", rest.DeleteClinic)

		r.Get("/clinic/{cid}/location", rest.GetAllLocation)
		r.Get("/clinic/location/{lid}", rest.GetLocation)
		r.Post("/clinic/location", rest.CreateLocation)
		r.Patch("/clinic/location/{lid}", rest.UpdateLocation)
		r.Delete("/clinic/location/{lid}", rest.DeleteLocation)
	})
}

func (rest *REST) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.Object[any]{Data: nil, Message: "OK"})
}

func (rest *REST) CreateClinic(w http.ResponseWriter, r *http.Request) {
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

	var req dto.RequestCreateClinic
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

	data, err := rest.clinicService.CreateClinic(ctx, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Create Clinic"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseCreateClinic]{Data: &data, Message: "OK"})
}

func (rest *REST) UpdateClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

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

	var req dto.RequestUpdateClinic
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

	data, err := rest.clinicService.UpdateClinic(ctx, id, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Update Clinic"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseUpdateClinic]{Data: &data, Message: "OK"})
}

func (rest *REST) DeleteClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	if err := rest.clinicService.DeleteClinic(ctx, id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Delete Clinic"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseUpdateClinic]{Message: "OK"})
}

func (rest *REST) GetClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	claims := ctx.Value(oauth.ClaimsContext)
	c, _ := json.Marshal(claims)
	var fClaims dto.FirebaseClaims
	json.Unmarshal(c, &fClaims)

	data, err := rest.clinicService.GetClinic(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Clinic"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseGetClinic]{Data: &data, Message: "OK"})
}

func (rest *REST) GetAllClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)

	claims := ctx.Value(oauth.ClaimsContext)
	c, _ := json.Marshal(claims)
	var fClaims dto.FirebaseClaims
	json.Unmarshal(c, &fClaims)

	data, err := rest.clinicService.GetAllClinic(ctx, dto.FilterGetClinic{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Clinic"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[[]dto.ResponseGetClinic]{Data: &data, Message: "OK"})
}

func (rest *REST) CreateLocation(w http.ResponseWriter, r *http.Request) {
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

	var req dto.RequestCreateLocation
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

	data, err := rest.clinicService.CreateLocation(ctx, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Create Location"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseCreateLocation]{Data: &data, Message: "OK"})
}

func (rest *REST) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "lid")

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

	var req dto.RequestUpdateLocation
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

	data, err := rest.clinicService.UpdateLocation(ctx, id, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Update Location"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseUpdateLocation]{Data: &data, Message: "OK"})
}

func (rest *REST) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "lid")

	if err := rest.clinicService.DeleteLocation(ctx, id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Delete Location"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseUpdateClinic]{Message: "OK"})

}

func (rest *REST) GetLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "lid")

	claims := ctx.Value(oauth.ClaimsContext)
	c, _ := json.Marshal(claims)
	var fClaims dto.FirebaseClaims
	json.Unmarshal(c, &fClaims)

	data, err := rest.clinicService.GetLocation(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Location"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseGetLocation]{Data: &data, Message: "OK"})
}

func (rest *REST) GetAllLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	cid := chi.URLParam(r, "cid")

	claims := ctx.Value(oauth.ClaimsContext)
	c, _ := json.Marshal(claims)
	var fClaims dto.FirebaseClaims
	json.Unmarshal(c, &fClaims)

	data, err := rest.clinicService.GetLocationByClinic(ctx, cid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Clinic"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[[]dto.ResponseGetLocation]{Data: &data, Message: "OK"})
}
