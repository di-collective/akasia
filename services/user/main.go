package main

import (
	"context"
	"fmt"
	"monorepo/internal/config"
	"monorepo/internal/db"
	"monorepo/internal/repository"
	"monorepo/services/user/api"
	"monorepo/services/user/models"
	"monorepo/services/user/service"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"
)

func main() {
	godotenv.Load()
	cfg := &config.Environment{}
	if err := env.Parse(cfg); err != nil {
		logrus.Fatalf("Failed to parse environment variables: %v", err)
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile(cfg.FirebaseConfig)
	fba, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logrus.Fatalf("Failed to initialize firebase app: %v", err)
	}

	fbaClient, err := fba.Auth(ctx)
	if err != nil {
		logrus.Fatalf("Failed to initialize firebase auth: %v", err)
	}

	pgdb := db.MustConnectPostgres(&db.PostgresConfig{
		SSLMode: cfg.DbSslMode,
		Name:    cfg.DbName,
		Host:    cfg.DbHost,
		Port:    cfg.DbPort,
		User:    cfg.DbUser,
		Pass:    cfg.DbPass,
	})

	mailer := gomail.NewMessage()
	dialer := gomail.NewDialer(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPAuthEmail,
		cfg.SMTPAuthPassword,
	)

	tbUser := repository.NewRepository[models.User, string](pgdb, repository.Tables.User)
	tbProfile := repository.NewRepository[models.Profile, string](pgdb, repository.Tables.Profile)

	restAPI := api.NewREST(
		service.NewOauthVerifier(tbUser, fbaClient, cfg),
		service.NewUserService(tbUser, tbProfile),
		service.NewEmailService(dialer, mailer, tbUser, tbProfile),
		cfg,
	)

	restAPI.InitializeRoutes()
	servicePort := fmt.Sprintf(":%d", cfg.ServicePort)
	logrus.Infof("Starting HTTP server in port: %d", cfg.ServicePort)
	if err := http.ListenAndServe(servicePort, restAPI.Router); err != nil {
		logrus.Fatalf("Failed to start HTTP server: %v", err)
	}
}
