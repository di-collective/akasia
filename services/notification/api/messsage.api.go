package api

import (
	"encoding/json"
	"errors"
	"monorepo/internal/dto"
	"monorepo/pkg/utils"
	"net/http"

	service "monorepo/services/notification/service"

	"github.com/go-chi/chi/v5"
)

func (rest *REST) ListMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Decode URL Query into DTO
	query := &dto.RequestListNotificationMessage{}
	err := rest.decoder.Decode(query, r.URL.Query())
	if err != nil {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "Query should not be empty"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 2. Pass DTO to be processed in the service layer
	messages, err := rest.service.ListMessages(ctx, query)
	if err != nil {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Convert service layer result into response DTO
	// 4. Write response
	data := utils.Map(messages, mapToNotificationMessage)
	json.NewEncoder(w).Encode(dto.Object[[]dto.NotificationMessage]{Data: &data})
	w.WriteHeader(http.StatusOK)
}

func (rest *REST) CreateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := &dto.RequestMutateNotificationMessage{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "Body should not be empty"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, err := rest.service.CreateMessage(ctx, body)
	if err != nil {
		isValidationError := false
		response := dto.Object[any]{Error: err.Error()}
		response.Error, isValidationError = mapValidationError(err)

		json.NewEncoder(w).Encode(response)
		w.WriteHeader(utils.Ternary(
			isValidationError,
			http.StatusBadRequest, http.StatusInternalServerError,
		))
		return
	}

	data := mapToNotificationMessage(message, 0)
	json.NewEncoder(w).Encode(dto.Object[dto.NotificationMessage]{
		Data: &data,
	})
	w.WriteHeader(http.StatusCreated)
}

func (rest *REST) GetMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Get URL parameter
	id := chi.URLParam(r, "id")

	// 2. Pass message id to be processed in the service layer
	message, err := rest.service.GetMessage(ctx, id)
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
	data := mapToNotificationMessage(message, 0)
	json.NewEncoder(w).Encode(dto.Object[dto.NotificationMessage]{Data: &data})
	w.WriteHeader(http.StatusOK)
}

func (rest *REST) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	body := &dto.RequestMutateNotificationMessage{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "Body should not be empty"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, err := rest.service.UpdateMessage(ctx, id, body)
	if err != nil {
		isValidationError := false
		response := dto.Object[any]{Error: err.Error()}
		response.Error, isValidationError = mapValidationError(err)

		json.NewEncoder(w).Encode(response)
		w.WriteHeader(utils.Ternary(
			isValidationError,
			http.StatusBadRequest, http.StatusInternalServerError,
		))
		return
	}

	data := mapToNotificationMessage(message, 0)
	json.NewEncoder(w).Encode(dto.Object[dto.NotificationMessage]{
		Data: &data,
	})
	w.WriteHeader(http.StatusOK)
}

func (rest *REST) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	w.Write(nil)
	w.WriteHeader(200)
}
