package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/models"
	"regexp"
	"time"
)

func (svc *Service) SaveLog(ctx context.Context, fluentbitLog *models.FluentBitReq) error {

	Logger := logger.CreateFileLoggerWithCtx(ctx)

	log, err := processFluentBitLog(fluentbitLog)
	if err != nil {
		Logger.Errorf("Failed to process fluentbit log: %v", err)
		return err
	}

	err = svc.repo.SaveLog(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) SaveBulkLog(ctx context.Context, fluentBitLogs []*models.FluentBitReq) error {

	Logger := logger.CreateFileLoggerWithCtx(ctx)

	var logs []*models.Log

	for _, fluentbitLog := range fluentBitLogs {
		log, err := processFluentBitLog(fluentbitLog)
		if err != nil {
			Logger.Errorf("Failed to process fluentbit log: %v", err)
			return err
		}

		logs = append(logs, log)
	}

	err := svc.repo.SaveBulkLog(ctx, logs)
	if err != nil {
		return err
	}
	return nil
}

func processFluentBitLog(fluentbitLog *models.FluentBitReq) (*models.Log, error) {
	parsedTime, err := time.Parse(time.RFC3339Nano, fluentbitLog.LogDecoded.Timestamp)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02T15:04:05.000-0700", fluentbitLog.LogDecoded.Timestamp) // your format
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %w", err)
		}
	}

	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleanLevel := re.ReplaceAllString(fluentbitLog.LogDecoded.Level, "")

	log := &models.Log{
		Timestamp: parsedTime,
		Namespace: fluentbitLog.LogDecoded.Namespace,
		Host:      fluentbitLog.LogDecoded.Host,
		Service:   fluentbitLog.LogDecoded.Service,
		Level:     cleanLevel,
		TraceID:   &fluentbitLog.LogDecoded.TraceID,
		Source:    fluentbitLog.Log,
	}

	var logMap map[string]interface{}
	data, _ := json.Marshal(fluentbitLog.LogDecoded)
	json.Unmarshal(data, &logMap)

	for k, v := range logMap {
		switch val := v.(type) {
		case string:
			log.StringNames = append(log.StringNames, k)
			log.StringValues = append(log.StringValues, val)
		case int:
			log.IntNames = append(log.IntNames, k)
			log.IntValues = append(log.IntValues, int64(val))
		case int64:
			log.IntNames = append(log.IntNames, k)
			log.IntValues = append(log.IntValues, val)
		case float64:
			log.FloatNames = append(log.FloatNames, k)
			log.FloatValues = append(log.FloatValues, val)
		case bool:
			log.BoolNames = append(log.BoolNames, k)
			log.BoolValues = append(log.BoolValues, fmt.Sprintf("%v", val))
		default:
			// Skip unsupported types
			return nil, fmt.Errorf("unsupported type for key %s: %T", k, val)
		}
	}

	return log, nil
}
