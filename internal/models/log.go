package models

import (
	"time"
)

type Log struct {
	Timestamp    time.Time              `json:"timestamp" db:"_timestamp"`
	Namespace    string                 `json:"namespace" db:"_namespace"`
	Host         string                 `json:"host" db:"host"`
	Service      string                 `json:"service" db:"service"`
	Level        string                 `json:"level" db:"level"`
	UserID       *uint64                `json:"user_id,omitempty" db:"user_id"`
	SessionID    *string                `json:"session_id,omitempty" db:"session_id"`
	TraceID      *string                `json:"trace_id,omitempty" db:"trace_id"`
	Data         map[string]interface{} `json:"data"` // dynamic fields
	Source       string                 `json:"_source" db:"_source"`
	StringNames  []string               `json:"string_names" db:"string_names"`
	StringValues []string               `json:"string_values" db:"string_values"`
	IntNames     []string               `json:"int_names" db:"int_names"`
	IntValues    []int64                `json:"int_values" db:"int_values"`
	FloatNames   []string               `json:"float_names" db:"float_names"`
	FloatValues  []float64              `json:"float_values" db:"float_values"`
	BoolNames    []string               `json:"bool_names" db:"bool_names"`
	BoolValues   []string               `json:"bool_values" db:"bool_values"`
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

type FluentBitReq struct {
	Date       float64 `json:"date"`
	Log        string  `json:"log"`
	LogDecoded RawLog  `json:"log_decoded"`
}
