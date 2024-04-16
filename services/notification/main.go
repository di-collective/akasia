package main

import (
	"fmt"
	"monorepo/internal/config"
	"monorepo/internal/db"
	"monorepo/services/notification/api"
	"monorepo/services/notification/models"
	"monorepo/services/notification/service"
	"net/http"

	"monorepo/internal/repository"

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

	tbMessage := repository.NewRepository[models.Message, string](pgdb, repository.Tables.Message)
	tbUserMessage := repository.NewRepository[models.UserMessage, string](pgdb, repository.Tables.UserMessage)
	vwUserMessage := repository.NewRepository[models.ViewUserMessage, string](pgdb, repository.Views.UserMessage)
	notifService := service.NewNotificationService(tbMessage, tbUserMessage, vwUserMessage)
	restAPI := api.NewREST(notifService)

	restAPI.InitializeRoutes()

	servicePort := fmt.Sprintf(":%d", cfg.ServicePort)
	logrus.Infof("Starting HTTP server in port: %d", cfg.ServicePort)
	if err := http.ListenAndServe(servicePort, restAPI.Router); err != nil {
		logrus.Fatalf("Failed to start HTTP server: %v", err)
	}
}
