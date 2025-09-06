package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type HistoryAction string

const (
	ActionCreate HistoryAction = "create"
	ActionUpdate HistoryAction = "update"
	ActionDelete HistoryAction = "delete"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), &j)
}

type ItemHistory struct {
	ID        int           `json:"id" db:"id"`
	ItemID    int           `json:"item_id" db:"item_id"`
	Action    HistoryAction `json:"action" db:"action"`
	ChangedBy string        `json:"changed_by" db:"changed_by"`
	OldValues JSONB         `json:"old_values" db:"old_values"`
	NewValues JSONB         `json:"new_values" db:"new_values"`
	ChangedAt time.Time     `json:"changed_at" db:"changed_at"`
}
