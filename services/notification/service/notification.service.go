package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"monorepo/internal/dto"
	"monorepo/pkg/common"

	"monorepo/internal/repository"
	"monorepo/services/notification/models"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

func NewNotificationService(
	tbMessage common.Repository[models.Message, string],
	tbUserMessage common.Repository[models.UserMessage, string],
	vwUserMessage common.Repository[models.ViewUserMessage, string],
) *NotificationService {
	service := &NotificationService{}
	service.validate = validator.New()
	service.tables.message = tbMessage
	service.tables.userMessage = tbUserMessage
	service.views.userMessage = vwUserMessage

	return service
}

type NotificationService struct {
	validate *validator.Validate
	tables   struct {
		message     common.Repository[models.Message, string]
		userMessage common.Repository[models.UserMessage, string]
	}
	views struct {
		userMessage common.Repository[models.ViewUserMessage, string]
	}
}

func (service *NotificationService) ListMessages(ctx context.Context, query *dto.RequestListNotificationMessage) ([]*models.Message, error) {
	// 1. Validate query
	err := service.validate.StructCtx(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. Map query parameters into expression
	filter := []exp.Expression{}
	if query.IDs != nil && len(query.IDs) > 0 {
		filter = append(filter, goqu.C("id").In(query.IDs))
	}
	if query.Type != nil && len(query.Type) > 0 {
		filter = append(filter, goqu.C("type").In(query.IDs))
	}
	if strings.TrimSpace(query.Content) != "" {
		filter = append(filter, goqu.C("content").ILike(query.Content))
	}

	// 3. Pass over to repository layer
	messages, err := service.tables.message.List(ctx, &common.FilterOptions{
		Filter: filter,
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	return messages, nil
}

func (service *NotificationService) CreateMessage(ctx context.Context, body *dto.RequestMutateNotificationMessage) (*models.Message, error) {
	err := service.validate.Struct(body)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrValidationFailed, err)
	}

	entity := &models.Message{
		ID:        ulid.Make().String(),
		CreatedAt: time.Now(),
		Type:      body.Type,
		Criteria:  sql.NullString{String: body.Criteria, Valid: true},
		Content:   body.Content,
	}
	if body.ScheduledAt != nil {
		entity.ScheduledAt.Time = *body.ScheduledAt
	}

	err = service.tables.message.Create(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	return entity, nil
}

func (service *NotificationService) GetMessage(ctx context.Context, id string) (*models.Message, error) {
	entity, err := service.tables.message.Get(ctx, id)
	switch {
	case err != nil && errors.Is(err, repository.ErrNoResult):
		return nil, fmt.Errorf("%w; %w", ErrNoResult, err)
	case err != nil:
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	return entity, nil
}

func (service *NotificationService) UpdateMessage(ctx context.Context, id string, body *dto.RequestMutateNotificationMessage) (*models.Message, error) {
	err := service.validate.Struct(body)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrValidationFailed, err)
	}

	entity := &models.Message{
		ID:       id,
		Type:     body.Type,
		Criteria: sql.NullString{String: body.Criteria, Valid: true},
		Content:  body.Content,
	}
	if body.ScheduledAt != nil {
		entity.ScheduledAt.Time = *body.ScheduledAt
	}

	err = service.tables.message.Update(ctx, id, entity)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	return service.GetMessage(ctx, id)
}

func (service *NotificationService) ListUserMessages(ctx context.Context, userId string) ([]*models.ViewUserMessage, error) {
	messages, err := service.views.userMessage.List(ctx, &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Select: []any{"message_id", "message_content", "read_at"},
		Filter: []exp.Expression{
			goqu.Or(
				goqu.C("message_type").Eq("all"),
				goqu.And(
					goqu.C("message_type").Eq("individual"),
					goqu.C("message_criteria").Eq(userId),
				),
			),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	return messages, nil
}

func (service *NotificationService) GetUserMessage(ctx context.Context, userId, messageId string) (*models.ViewUserMessage, error) {
	messages, err := service.views.userMessage.List(ctx, &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Select: []any{"message_id", "message_content", "read_at"},
		Filter: []exp.Expression{
			goqu.C("message_id").Eq(messageId),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	if len(messages) <= 0 {
		return nil, fmt.Errorf("%w; %w", ErrNoResult, err)
	}

	return messages[0], nil
}

func (service *NotificationService) ReadUserMessage(ctx context.Context, userId, messageId string) error {
	messages, err := service.tables.userMessage.List(ctx, &common.FilterOptions{
		Filter: []exp.Expression{
			goqu.C("user_id").Eq(userId),
			goqu.C("message_id").Eq(messageId),
		},
	})
	if err != nil {
		return fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	}

	now := time.Now()
	if len(messages) <= 0 {
		return service.tables.userMessage.Create(ctx, &models.UserMessage{
			ID:        ulid.Make().String(),
			UserID:    userId,
			MessageID: messageId,
			CreatedAt: now,
			ReadAt:    sql.NullTime{Time: now, Valid: true},
		})
	}

	usm := messages[0]
	usm.ReadAt.Time = now
	usm.ReadAt.Valid = true
	return service.tables.userMessage.Update(ctx, usm.ID, usm)
}
