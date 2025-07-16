package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/models"
)

func (repo *Repo) SaveLog(ctx context.Context, log *models.Log) error {

	Logger := logger.CreateFileLoggerWithCtx(ctx)
	var (
		stringNames  []string
		stringValues []string
		intNames     []string
		intValues    []int64
		floatNames   []string
		floatValues  []float64
		boolNames    []string
		boolValues   []string
	)

	// Flatten dynamic map
	for k, v := range log.Data {
		switch val := v.(type) {
		case string:
			stringNames = append(stringNames, k)
			stringValues = append(stringValues, val)
		case int:
			intNames = append(intNames, k)
			intValues = append(intValues, int64(val))
		case int64:
			intNames = append(intNames, k)
			intValues = append(intValues, val)
		case float64:
			floatNames = append(floatNames, k)
			floatValues = append(floatValues, val)
		case bool:
			boolNames = append(boolNames, k)
			boolValues = append(boolValues, fmt.Sprintf("%v", val))
		default:
			// Skip unsupported types
			Logger.Errorf("Unsupported type for key %s: %T", k, v)
		}
	}

	// Store original data as JSON
	sourceJSON, err := json.Marshal(log.Data)
	if err != nil {
		return err
	}

	// Insert into ClickHouse
	query := `
		INSERT INTO logs (
			_timestamp, _namespace, host, service, level,
			user_id, session_id, trace_id,
			_source,
			string_names, string_values,
			int_names, int_values,
			float_names, float_values,
			bool_names, bool_values
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	err = repo.DB.DB().Exec(ctx, query,
		log.Timestamp,
		log.Namespace,
		log.Host,
		log.Service,
		log.Level,
		log.UserID,
		log.SessionID,
		log.TraceID,
		string(sourceJSON),
		stringNames, stringValues,
		intNames, intValues,
		floatNames, floatValues,
		boolNames, boolValues,
	)

	if err != nil {
		Logger.Errorf("Error inserting log: %v", err)
		return err
	}
	return nil
}
