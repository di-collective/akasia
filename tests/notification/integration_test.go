package tests

import (
	"context"

	"monorepo/internal/config"
	"monorepo/internal/db"
	"monorepo/pkg/common"

	"monorepo/internal/repository"
	"monorepo/services/notification/models"
	"os"
	"testing"
	"time"

	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	pgdb    *sqlx.DB
	msgRepo *repository.Repository[models.Message, string]
)

func TestMain(m *testing.M) {
	godotenv.Load("test.env")
	cfg := &config.Environment{}
	if err := env.Parse(cfg); err != nil {
		logrus.Fatalf("Failed to parse environment variables: %v", err)
	}

	pgdb = db.MustConnectPostgres(&db.PostgresConfig{
		SSLMode: cfg.DbSslMode,
		Name:    cfg.DbName,
		Host:    cfg.DbHost,
		Port:    cfg.DbPort,
		User:    cfg.DbUser,
		Pass:    cfg.DbPass,
	})

	msgRepo = repository.NewRepository[models.Message, string](pgdb, repository.Tables.Message)
	os.Exit(m.Run())
}

func TestConnection(t *testing.T) {
	assert.NotNil(t, pgdb)
	assert.NotNil(t, msgRepo)
}

func TestMessageRepository(t *testing.T) {
	ctx := context.Background()

	_, err := msgRepo.Raw(ctx, "truncate table public.message")
	assert.NoError(t, err)

	// 1. Create
	id := ulid.Make()
	err = msgRepo.Create(ctx, &models.Message{
		ID:        id.String(),
		CreatedAt: time.Now(),
		Type:      "all",
		Content:   "random",
	})
	assert.NoError(t, err)

	// 2. Get
	msg, err := msgRepo.Get(ctx, id.String())
	assert.NoError(t, err)
	assert.Equal(t, id.String(), msg.ID)

	// 3. Delete
	err = msgRepo.Delete(ctx, id.String())
	assert.NoError(t, err)

	// 4. Get non-existing (already deleted)
	msg, err = msgRepo.Get(ctx, id.String())
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNoResult)
	assert.Nil(t, msg)

	// 5. Create 5 messages
	var idBasket = [5]string{}
	for i := 0; i < 5; i++ {
		id := ulid.Make()
		idBasket[i] = id.String()
		err := msgRepo.Create(ctx, &models.Message{
			ID:        id.String(),
			CreatedAt: time.Now(),
			Type:      "all",
			Content:   "random",
		})
		assert.NoError(t, err)
	}

	// 6. Find all messages
	messages, err := msgRepo.List(ctx, &common.FilterOptions{})
	assert.NoError(t, err)
	assert.Len(t, messages, 5)

	// 7. Delete
	for i := 0; i < 5; i++ {
		err = msgRepo.Delete(ctx, idBasket[i])
		assert.NoError(t, err)
	}
}
