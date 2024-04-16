package api

import (
	"errors"
	"fmt"
	"monorepo/internal/dto"
	"monorepo/pkg/utils"
	"monorepo/services/notification/models"

	"github.com/go-playground/validator/v10"
)

func mapToNotificationMessage(entity *models.Message, _ int) dto.NotificationMessage {
	return dto.NotificationMessage{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		ScheduledAt: utils.Ternary(entity.ScheduledAt.Valid, &entity.ScheduledAt.Time, nil),
		SentAt:      utils.Ternary(entity.SentAt.Valid, &entity.SentAt.Time, nil),
		Type:        entity.Type,
		Criteria:    entity.Criteria.String,
		Content:     entity.Content,
	}
}

func mapToNotificationUserMessage(entity *models.ViewUserMessage, _ int) dto.NotificationUserMessage {
	return dto.NotificationUserMessage{
		ID:      entity.ID,
		Content: entity.Content,
		ReadAt:  utils.Ternary(entity.ReadAt.Valid, &entity.ReadAt.Time, nil),
	}
}

// returns true if err is of type ValidationErrors
func mapValidationError(err error) ([]string, bool) {
	verr := validator.ValidationErrors{}
	isValidationError := errors.As(err, &verr)
	if isValidationError {
		return utils.Map(verr, func(ferr validator.FieldError, i int) string {
			return utils.Ternary(
				ferr.Param() == "",
				fmt.Sprintf(`Field: '%s', failed on: '%s' spec`, ferr.Field(), ferr.Tag()),
				fmt.Sprintf(`Field: '%s', failed on: '%s:%s' spec`, ferr.Field(), ferr.Tag(), ferr.Param()),
			)
		}), true
	}

	return nil, false
}
