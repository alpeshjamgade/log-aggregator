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
	LoginID      string                 `json:"login_id,omitempty" db:"login_id"`
	ClientID     string                 `json:"client_id,omitempty" db:"client_id"`
	SessionID    string                 `json:"session_id,omitempty" db:"session_id"`
	TraceID      string                 `json:"trace_id,omitempty" db:"trace_id"`
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
	TS        string `json:"ts"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	TraceID   string `json:"trace_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Caller    string `json:"caller,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Host      string `json:"host,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	LoginID   string `json:"login_id,omitempty"`
}

type FluentBitReq struct {
	Date       float64            `json:"date"`
	Log        string             `json:"log"`
	LogDecoded RawLog             `json:"log_decoded"`
	Kubernetes KubernetesMetadata `json:"kubernetes"`
}

type KubernetesMetadata struct {
	PodName        string          `json:"pod_name"`
	NamespaceName  string          `json:"namespace_name"`
	PodID          string          `json:"pod_id"`
	Labels         KubernetesLabel `json:"labels"`
	Annotations    map[string]any  `json:"annotations"`
	Host           string          `json:"host"`
	ContainerName  string          `json:"container_name"`
	DockerID       string          `json:"docker_id"`
	ContainerHash  string          `json:"container_hash"`
	ContainerImage string          `json:"container_image"`
}

type KubernetesLabel struct {
	App string `json:"app"`
}
