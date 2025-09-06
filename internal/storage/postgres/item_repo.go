package postgres

import (
	"WarehouseControl/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type ItemStorageI interface {
	CreateItem(ctx context.Context, item *models.Item, changedBy string) error
	GetAllItems(ctx context.Context) ([]*models.Item, error)
	GetItemByID(ctx context.Context, id int) (*models.Item, error)
	UpdateItem(ctx context.Context, item *models.Item, changedBy string) error
	DeleteItem(ctx context.Context, id int, changedBy string) error
}

type ItemStorage struct {
	db *sql.DB
}

func NewItemStorage(db *sql.DB) *ItemStorage {
	return &ItemStorage{db: db}
}

func (s *ItemStorage) CreateItem(ctx context.Context, item *models.Item, changedBy string) error {
	_, err := s.db.ExecContext(ctx, "SET LOCAL app.user = $1", changedBy)
	if err != nil {
		return fmt.Errorf("failed to set user context: %w", err)
	}

	query := `INSERT INTO items (name, quantity) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err = s.db.QueryRowContext(ctx, query, item.Name, item.Quantity).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}
	return nil
}

func (s *ItemStorage) GetAllItems(ctx context.Context) ([]*models.Item, error) {
	query := `SELECT id, name, quantity, created_at, updated_at FROM items ORDER BY id`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

func (s *ItemStorage) GetItemByID(ctx context.Context, id int) (*models.Item, error) {
	query := `SELECT id, name, quantity, created_at, updated_at FROM items WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)

	var item models.Item
	err := row.Scan(&item.ID, &item.Name, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &item, nil
}

func (s *ItemStorage) UpdateItem(ctx context.Context, item *models.Item, changedBy string) error {
\	_, err := s.db.ExecContext(ctx, "SET LOCAL app.user = $1", changedBy)
	if err != nil {
		return fmt.Errorf("failed to set user context: %w", err)
	}

	query := `UPDATE items SET name = $1, quantity = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`
	err = s.db.QueryRowContext(ctx, query, item.Name, item.Quantity, item.ID).Scan(&item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("item not found")
		}
		return fmt.Errorf("failed to update item: %w", err)
	}
	return nil
}

func (s *ItemStorage) DeleteItem(ctx context.Context, id int, changedBy string) error {
	_, err := s.db.ExecContext(ctx, "SET LOCAL app.user = $1", changedBy)
	if err != nil {
		return fmt.Errorf("failed to set user context: %w", err)
	}

	query := `DELETE FROM items WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}