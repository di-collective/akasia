package api

import (
	"encoding/json"
	"errors"
	"monorepo/internal/dto"
	"monorepo/pkg/utils"
	service "monorepo/services/notification/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (rest *REST) ListUserMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := chi.URLParam(r, "userId")

	messages, err := rest.service.ListUserMessages(ctx, userId)
	if err != nil {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Convert service layer result into response DTO
	data := utils.Map(messages, mapToNotificationUserMessage)
	json.NewEncoder(w).Encode(dto.Object[[]dto.NotificationUserMessage]{Data: &data})
	w.WriteHeader(http.StatusOK)
}

func (rest *REST) GetUserMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := chi.URLParam(r, "userId")
	messageId := chi.URLParam(r, "messageId")

	message, err := rest.service.GetUserMessage(ctx, userId, messageId)
	switch {
	case err != nil && errors.Is(err, service.ErrNoResult):
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Convert service layer result into response DTO
	data := mapToNotificationUserMessage(message, 0)
	json.NewEncoder(w).Encode(dto.Object[dto.NotificationUserMessage]{Data: &data})
	w.WriteHeader(http.StatusOK)
}

func (rest *REST) ReadUserMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := chi.URLParam(r, "userId")
	messageId := chi.URLParam(r, "messageId")

	err := rest.service.ReadUserMessage(ctx, userId, messageId)
	if err != nil {

		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(nil)
	w.WriteHeader(200)
}
