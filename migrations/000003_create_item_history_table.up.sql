-- +goose Up
CREATE TABLE item_history
(
    id         SERIAL PRIMARY KEY,
    item_id    INTEGER     NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    action     VARCHAR(20) NOT NULL, -- 'create', 'update', 'delete'
    changed_by VARCHAR(50),          -- имя пользователя
    old_values JSONB,
    new_values JSONB,
    changed_at TIMESTAMP DEFAULT NOW()
);