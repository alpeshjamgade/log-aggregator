package models

import (
	"time"
)

type Log struct {
	Timestamp time.Time              `json:"timestamp" db:"_timestamp"`
	Namespace string                 `json:"namespace" db:"_namespace"`
	Host      string                 `json:"host" db:"host"`
	Service   string                 `json:"service" db:"service"`
	Level     string                 `json:"level" db:"level"`
	UserID    *uint64                `json:"user_id,omitempty" db:"user_id"`
	SessionID *string                `json:"session_id,omitempty" db:"session_id"`
	TraceID   *string                `json:"trace_id,omitempty" db:"trace_id"`
	Data      map[string]interface{} `json:"data"` // dynamic fields
}

type RawLog struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	TraceID   string `json:"trace_id,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Caller    string `json:"caller,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Host      string `json:"host,omitempty"`
}
