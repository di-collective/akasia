package main

import (
	"fmt"
	"monorepo/internal/config"
	"monorepo/internal/db"
	"monorepo/internal/repository"
	"monorepo/services/calendar/api"
	"monorepo/services/calendar/models"
	"monorepo/services/calendar/service"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	cfg := &config.Environment{}
	if err := env.Parse(cfg); err != nil {
		logrus.Fatalf("Failed to parse environment variables: %v", err)
	}

	pgdb := db.MustConnectPostgres(&db.PostgresConfig{
		SSLMode: cfg.DbSslMode,
		Name:    cfg.DbName,
		Host:    cfg.DbHost,
		Port:    cfg.DbPort,
		User:    cfg.DbUser,
		Pass:    cfg.DbPass,
	})

	tbEvent := repository.NewRepository[models.Event, string](pgdb, repository.Tables.Event)

	eventService := service.NewEventService(tbEvent)
	restAPI := api.NewREST(eventService, cfg)

	restAPI.InitializeRoutes()

	servicePort := fmt.Sprintf(":%d", cfg.ServicePort)
	logrus.Infof("Starting HTTP server in port: %d", cfg.ServicePort)
	if err := http.ListenAndServe(servicePort, restAPI.Router); err != nil {
		logrus.Fatalf("Failed to start HTTP server: %v", err)
	}
}
