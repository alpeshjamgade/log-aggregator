package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log-aggregator/internal/constants"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/models"
	"regexp"
	"time"
)

func (svc *Service) SaveLog(ctx context.Context, fluentbitLog *models.FluentBitReq) error {

	Logger := logger.CreateFileLoggerWithCtx(ctx)

	log, err := processFluentBitLog(ctx, fluentbitLog)
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
		log, err := processFluentBitLog(ctx, fluentbitLog)
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

func processFluentBitLog(ctx context.Context, fluentbitLog *models.FluentBitReq) (*models.Log, error) {

	log := &models.Log{
		Namespace: fluentbitLog.LogDecoded.Namespace,
		Host:      fluentbitLog.Kubernetes.Host,
		Service:   fluentbitLog.Kubernetes.Labels.App,
		TraceID:   fluentbitLog.LogDecoded.TraceID,
		UserID:    fluentbitLog.LogDecoded.LoginID,
		Source:    fluentbitLog.Log,
	}

	err := setTimestamp(ctx, log, fluentbitLog)
	if err != nil {
		return nil, err
	}

	setTraceID(ctx, log, fluentbitLog)

	err = setFieldNameWithValues(ctx, log, fluentbitLog)
	if err != nil {
		return nil, err
	}

	setLogLevel(ctx, log, fluentbitLog)

	return log, nil
}

func setTimestamp(ctx context.Context, log *models.Log, fluentbitLog *models.FluentBitReq) error {
	Logger := logger.CreateFileLoggerWithCtx(ctx)
	layouts := []string{
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
	}

	var timestamp string
	if fluentbitLog.LogDecoded.Timestamp != "" {
		timestamp = fluentbitLog.LogDecoded.Timestamp
	}

	if timestamp == constants.Empty && fluentbitLog.LogDecoded.TS != "" {
		timestamp = fluentbitLog.LogDecoded.TS
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		for _, layout := range layouts {
			parsedTime, err = time.Parse(layout, timestamp)
			if err == nil {
				break
			}
		}
		if err != nil {
			Logger.Errorw("invalid timestamp", "timestamp", timestamp, "error", err)
			return err
		}
	}

	log.Timestamp = parsedTime

	return nil
}

func setTraceID(ctx context.Context, log *models.Log, fluentbitLog *models.FluentBitReq) {
	if fluentbitLog.LogDecoded.TraceID != "" {
		log.TraceID = fluentbitLog.LogDecoded.TraceID
	}

	// check for requestID if traceID is still empty
	if log.TraceID == "" && fluentbitLog.LogDecoded.RequestID != "" {
		log.TraceID = fluentbitLog.LogDecoded.RequestID
	}
}

func setFieldNameWithValues(ctx context.Context, log *models.Log, fluentbitLog *models.FluentBitReq) error {
	var logMap map[string]interface{}
	data, _ := json.Marshal(fluentbitLog.LogDecoded)
	err := json.Unmarshal(data, &logMap)
	if err != nil {
		return err
	}

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
			return fmt.Errorf("unsupported type for key %s: %T", k, val)
		}
	}

	return nil
}

func setLogLevel(ctx context.Context, log *models.Log, fluentbitLog *models.FluentBitReq) {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleanLevel := re.ReplaceAllString(fluentbitLog.LogDecoded.Level, "")

	log.Level = cleanLevel
}
