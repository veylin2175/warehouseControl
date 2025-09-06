package postgres

import (
	"WarehouseControl/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type HistoryStorageI interface {
	GetHistoryByItemID(ctx context.Context, itemID int) ([]*models.ItemHistory, error)
	GetAllHistory(ctx context.Context) ([]*models.ItemHistory, error)
}

type HistoryStorage struct {
	db *sql.DB
}

func NewHistoryStorage(db *sql.DB) *HistoryStorage {
	return &HistoryStorage{db: db}
}

func (s *HistoryStorage) GetHistoryByItemID(ctx context.Context, itemID int) ([]*models.ItemHistory, error) {
	query := `SELECT id, item_id, action, changed_by, old_values, new_values, changed_at 
	          FROM item_history WHERE item_id = $1 ORDER BY changed_at DESC`
	rows, err := s.db.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	defer rows.Close()

	var history []*models.ItemHistory
	for rows.Next() {
		var h models.ItemHistory
		err = rows.Scan(&h.ID, &h.ItemID, &h.Action, &h.ChangedBy, &h.OldValues, &h.NewValues, &h.ChangedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history record: %w", err)
		}
		history = append(history, &h)
	}

	return history, nil
}

func (s *HistoryStorage) GetAllHistory(ctx context.Context) ([]*models.ItemHistory, error) {
	query := `SELECT id, item_id, action, changed_by, old_values, new_values, changed_at 
	          FROM item_history ORDER BY changed_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all history: %w", err)
	}
	defer rows.Close()

	var history []*models.ItemHistory
	for rows.Next() {
		var h models.ItemHistory
		err = rows.Scan(&h.ID, &h.ItemID, &h.Action, &h.ChangedBy, &h.OldValues, &h.NewValues, &h.ChangedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history record: %w", err)
		}
		history = append(history, &h)
	}

	return history, nil
}
