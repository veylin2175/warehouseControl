package handlers

import (
	"WarehouseControl/internal/http-server/handlers/middleware"
	"WarehouseControl/internal/lib/api/response"
	"WarehouseControl/internal/models"
	"WarehouseControl/internal/storage/postgres"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemsHandler struct {
	itemStorage postgres.ItemStorageI
	log         *slog.Logger
}

func NewItemsHandler(itemStorage postgres.ItemStorageI, log *slog.Logger) *ItemsHandler {
	return &ItemsHandler{
		itemStorage: itemStorage,
		log:         log,
	}
}

type createItemRequest struct {
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity"`
}

func (h *ItemsHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.items.CreateItem"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Header.Get("X-Request-ID")),
	)

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Error("user not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("internal server error"))
		return
	}

	var req createItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.Error("invalid request body"))
		return
	}

	item := &models.Item{
		Name:     req.Name,
		Quantity: req.Quantity,
	}

	if err := h.itemStorage.CreateItem(r.Context(), item, claims.Username); err != nil {
		log.Error("failed to create item", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("failed to create item"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		response.Response
		Data *models.Item `json:"data,omitempty"`
	}{
		Response: response.Response{Status: response.StatusOK},
		Data:     item,
	})
}

func (h *ItemsHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.items.GetAllItems"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Header.Get("X-Request-ID")),
	)

	items, err := h.itemStorage.GetAllItems(r.Context())
	if err != nil {
		log.Error("failed to get items", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("failed to get items"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		response.Response
		Data []*models.Item `json:"data,omitempty"`
	}{
		Response: response.Response{Status: response.StatusOK},
		Data:     items,
	})
}

func (h *ItemsHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.items.GetItemByID"

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

	item, err := h.itemStorage.GetItemByID(r.Context(), id)
	if err != nil {
		log.Warn("item not found", slog.Int("id", id))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.Error("item not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		response.Response
		Data *models.Item `json:"data,omitempty"`
	}{
		Response: response.Response{Status: response.StatusOK},
		Data:     item,
	})
}

type updateItemRequest struct {
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity"`
}

func (h *ItemsHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.items.UpdateItem"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Header.Get("X-Request-ID")),
	)

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Error("user not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("internal server error"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid item id", slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.Error("invalid item id"))
		return
	}

	var req updateItemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.Error("invalid request body"))
		return
	}

	item := &models.Item{
		ID:       id,
		Name:     req.Name,
		Quantity: req.Quantity,
	}

	if err = h.itemStorage.UpdateItem(r.Context(), item, claims.Username); err != nil {
		log.Error("failed to update item", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("failed to update item"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		response.Response
		Data *models.Item `json:"data,omitempty"`
	}{
		Response: response.Response{Status: response.StatusOK},
		Data:     item,
	})
}

func (h *ItemsHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.items.DeleteItem"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Header.Get("X-Request-ID")),
	)

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Error("user not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("internal server error"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid item id", slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.Error("invalid item id"))
		return
	}

	if err = h.itemStorage.DeleteItem(r.Context(), id, claims.Username); err != nil {
		log.Error("failed to delete item", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("failed to delete item"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.OK())
}
