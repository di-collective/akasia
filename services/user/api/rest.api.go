package api

import (
	"encoding/json"
	"monorepo/internal/config"
	"monorepo/internal/dto"
	"monorepo/services/user/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"

	"github.com/gorilla/schema"
)

type REST struct {
	Router *chi.Mux

	decoder         *schema.Decoder
	userService     *service.UserService
	oauthServer     *oauth.BearerServer
	oauthVerifier   *service.OauthVerifier
	oauthAuthorizer func(next http.Handler) http.Handler
	env             *config.Environment
}

func NewREST(
	oauthVerifier *service.OauthVerifier,
	userService *service.UserService,
	env *config.Environment,
) *REST {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Compress(6))

	return &REST{
		Router:          r,
		decoder:         schema.NewDecoder(),
		userService:     userService,
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
	rest.Router.Group(func(r chi.Router) {
		r.Use(rest.oauthAuthorizer)

		r.Post("/me", rest.MyCredential)
	})
}

func (rest *REST) Healthcheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(dto.Object[any]{Data: nil, Message: "OK"})
	w.WriteHeader(200)
}

func (rest *REST) MyCredential(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := ctx.Value(oauth.ClaimsContext)
	json.NewEncoder(w).Encode(dto.Object[any]{Data: &claims, Message: "OK"})
	w.WriteHeader(200)
}
