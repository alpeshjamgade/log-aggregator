package service

import (
	"context"
	"fmt"
	"log-aggregator/internal/models"
	"regexp"
	"time"
)

func (svc *Service) SaveLog(ctx context.Context, raw *models.RawLog) error {
	parsedTime, err := time.Parse(time.RFC3339Nano, raw.Timestamp)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02T15:04:05.000-0700", raw.Timestamp) // your format
		if err != nil {
			return fmt.Errorf("invalid timestamp: %w", err)
		}
	}

	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleanLevel := re.ReplaceAllString(raw.Level, "")

	data := map[string]any{}

	if raw.Msg != "" {
		data["msg"] = raw.Msg
	}

	if raw.Caller != "" {
		data["caller"] = raw.Caller
	}

	log := &models.Log{
		Timestamp: parsedTime,
		Namespace: raw.Namespace,
		Host:      raw.Host,
		Service:   raw.Service,
		Level:     cleanLevel,
		TraceID:   ptr(raw.TraceID),
		Data:      data,
	}

	err = svc.repo.SaveLog(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func ptr[T any](v T) *T { return &v }
