package api

import (
	"encoding/json"
	"fmt"
	"io"
	"monorepo/internal/config"
	"monorepo/internal/dto"
	"monorepo/services/user/service"
	"net/http"
	"path/filepath"
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
	userService     *service.UserService
	emailService    *service.EmailService
	oauthServer     *oauth.BearerServer
	oauthVerifier   *service.OauthVerifier
	oauthAuthorizer func(next http.Handler) http.Handler
	env             *config.Environment
}

func NewREST(
	oauthVerifier *service.OauthVerifier,
	userService *service.UserService,
	emailService *service.EmailService,
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
		userService:     userService,
		emailService:    emailService,
		oauthServer:     oauth.NewBearerServer(env.JWTSecret, time.Hour*4, oauthVerifier, nil),
		oauthAuthorizer: oauth.Authorize(env.JWTSecret, nil),
		oauthVerifier:   oauthVerifier,
		env:             env,
	}
}

func (rest *REST) InitializeRoutes() {
	rest.Router.Get("/", rest.Healthcheck)
	rest.Router.Post("/credentials/login", rest.oauthServer.UserCredentials)
	rest.Router.Post("/credentials/firebase-auth", rest.FirebaseAuth)
	rest.Router.Post("/credentials/forgot-password", rest.ForgotPassword)
	rest.Router.Post("/credentials/update-password", rest.UpdatePassword)
	rest.Router.Group(func(r chi.Router) {
		r.Use(rest.oauthAuthorizer)

		r.Get("/me", rest.MyCredential)
		r.Get("/profile", rest.GetProfile)
		r.Post("/profile", rest.CreateProfile)
		r.Patch("/profile/{id}", rest.UpdateProfile)
		r.Patch("/profile/{id}/photo", rest.UploadPhoto)
		r.Delete("/profile/{id}", rest.DeleteProfile)
	})
}

func (rest *REST) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.Object[any]{Data: nil, Message: "OK"})
}

func (rest *REST) CreateProfile(w http.ResponseWriter, r *http.Request) {
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

	var req dto.RequestCreateProfile
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

	req.UserID = fClaims.UserID
	data, err := rest.userService.CreateProfile(ctx, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Create Profile"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[*dto.ResponseCreateProfile]{Data: &data, Message: "OK"})
}

func (rest *REST) MyCredential(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	claims := ctx.Value(oauth.ClaimsContext)
	json.NewEncoder(w).Encode(dto.Object[any]{Data: &claims, Message: "OK"})
}

func (rest *REST) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	claims := ctx.Value(oauth.ClaimsContext)

	fc := dto.FirebaseClaims{}
	c, err := json.Marshal(claims)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	err = json.Unmarshal(c, &fc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	data, err := rest.userService.GetProfile(ctx, &fc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Get Profile"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[any]{Data: &data, Message: "OK"})
}

func (rest *REST) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	var req dto.RequestUpdateProfile
	err = json.Unmarshal(payload, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		err = err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}

	data, err := rest.userService.UpdateProfile(ctx, userId, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Update Profile"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[any]{Data: &data, Message: "OK"})
}

func (rest *REST) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")

	err := rest.userService.DeleteProfile(ctx, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Delete Profile"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[any]{Message: "Profile deleted successfully"})
}

func (rest *REST) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		return
	}
	defer file.Close()

	t := time.Now()
	nanoT := t.Format("20060102150405")
	fileExt := filepath.Ext(handler.Filename)
	userId := chi.URLParam(r, "id")
	fileName := fmt.Sprintf("PP-%s-%s%s", userId, nanoT, fileExt)
	data, err := rest.userService.UploadPhoto(ctx, rest.env, file, fileName, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error(), Message: "Failed to Change Profile Picture"})
		return
	}

	json.NewEncoder(w).Encode(dto.Object[any]{Data: &data, Message: "OK"})
}
