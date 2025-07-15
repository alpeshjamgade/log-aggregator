package models

import (
	"time"
)

type Event struct {
	ID        string         `json:"id" db:"id"`
	Data      map[string]any `json:"data" db:"data"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
	LastLogin time.Time      `json:"last_login" db:"last_login"`
}
