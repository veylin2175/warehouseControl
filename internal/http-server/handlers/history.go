package handlers

import (
	"WarehouseControl/internal/lib/api/response"
	"WarehouseControl/internal/models"
	"WarehouseControl/internal/storage/postgres"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type HistoryHandler struct {
	historyStorage postgres.HistoryStorageI
	log            *slog.Logger
}

func NewHistoryHandler(historyStorage postgres.HistoryStorageI, log *slog.Logger) *HistoryHandler {
	return &HistoryHandler{
		historyStorage: historyStorage,
		log:            log,
	}
}

func (h *HistoryHandler) GetHistoryByItemID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.history.GetHistoryByItemID"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Header.Get("X-Request-ID")),
	)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid item id", slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.Error("invalid item id"))
		return
	}

	history, err := h.historyStorage.GetHistoryByItemID(r.Context(), id)
	if err != nil {
		log.Error("failed to get history", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("failed to get history"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		response.Response
		Data []*models.ItemHistory `json:"data,omitempty"`
	}{
		Response: response.Response{Status: response.StatusOK},
		Data:     history,
	})
}

func (h *HistoryHandler) GetAllHistory(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.history.GetAllHistory"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Header.Get("X-Request-ID")),
	)

	history, err := h.historyStorage.GetAllHistory(r.Context())
	if err != nil {
		log.Error("failed to get all history", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("failed to get history"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		response.Response
		Data []*models.ItemHistory `json:"data,omitempty"`
	}{
		Response: response.Response{Status: response.StatusOK},
		Data:     history,
	})
}
