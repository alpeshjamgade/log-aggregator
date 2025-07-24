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
			Logger.Warnf("Failed to process fluentbit. Skipping!! err: %v, log: %v", err, fluentbitLog)
			continue
		} else {
			logs = append(logs, log)
		}
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
		LoginID:   fluentbitLog.LogDecoded.LoginID,
		ClientID:  fluentbitLog.LogDecoded.ClientID,
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

	setUserID(ctx, log, fluentbitLog)

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

func setUserID(ctx context.Context, log *models.Log, fluentbitLog *models.FluentBitReq) {
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	if fluentbitLog.LogDecoded.LoginID != "" {
		log.LoginID = fluentbitLog.LogDecoded.LoginID
	}

	if fluentbitLog.LogDecoded.ClientID != "" {
		log.LoginID = fluentbitLog.LogDecoded.ClientID
	}

	clientID, err := ParseKeywordFromText(fluentbitLog.LogDecoded.Msg, "client_id")
	if err != nil {
		Logger.Debugf("error while parsing keyword %s from log message", "client_id")
	} else {
		log.ClientID = clientID
	}

	loginID, err := ParseKeywordFromText(fluentbitLog.LogDecoded.Msg, "login_id")
	if err != nil {
		Logger.Debugf("error while parsing keyword %s from log message", "login_id")
	} else {
		log.ClientID = loginID
	}
}

func ParseKeywordFromText(log string, key string) (string, error) {
	regexStr := fmt.Sprintf(`(?i)%s[^A-Za-z0-9]*([A-Za-z0-9_-]+)`, regexp.QuoteMeta(key))
	re := regexp.MustCompile(regexStr)
	matches := re.FindStringSubmatch(log)
	if len(matches) < 2 {
		return "", fmt.Errorf("key not found")
	}
	return matches[1], nil
}

func (svc *Service) SaveBulkLogV2(ctx context.Context, flattenLogs []map[string]any) error {
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	var logs []*models.Log

	for _, flattenLog := range flattenLogs {
		log, err := processFlattenLog(ctx, flattenLog)

		if err != nil {
			Logger.Warnf("Failed to process fluentbit. Skipping!! err: %v, log: %v", err, flattenLog)
			continue
		} else {
			logs = append(logs, log)
		}
	}

	err := svc.repo.SaveBulkLog(ctx, logs)
	if err != nil {

		Logger.Errorf("Error while saving log")
		return err
	}

	return nil
}

func processFlattenLog(ctx context.Context, flattenLog map[string]any) (*models.Log, error) {

	log := &models.Log{
		Host: flattenLog["kubernetes.host"].(string),
	}
	return log, nil
}
