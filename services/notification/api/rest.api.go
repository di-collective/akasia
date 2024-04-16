package api

import (
	"encoding/json"
	"monorepo/internal/dto"
	"monorepo/services/notification/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gorilla/schema"
)

type REST struct {
	Router  *chi.Mux
	service *service.NotificationService
	decoder *schema.Decoder
}

func NewREST(service *service.NotificationService) *REST {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Compress(6))

	return &REST{
		Router:  r,
		service: service,
		decoder: schema.NewDecoder(),
	}
}

func (rest *REST) InitializeRoutes() {
	rest.Router.Get("/healthcheck", rest.Healthcheck)

	rest.Router.Get("/messages", rest.ListMessages)
	rest.Router.Post("/messages", rest.CreateMessage)
	rest.Router.Get("/messages/{id}", rest.GetMessage)
	rest.Router.Put("/messages/{id}", rest.UpdateMessage)
	rest.Router.Delete("/messages/{id}", rest.DeleteMessage)

	rest.Router.Get("/users/{userId}/messages", rest.ListUserMessages)
	rest.Router.Get("/users/{userId}/messages/{messageId}", rest.GetUserMessage)
	rest.Router.Post("/users/{userId}/messages/{messageId}/read", rest.ReadUserMessage)
}

func (rest *REST) Healthcheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(dto.Object[any]{Data: nil, Message: "OK"})
	w.WriteHeader(200)
}
