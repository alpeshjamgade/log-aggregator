package repo

import (
	"context"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/models"
)

func (repo *Repo) SaveLog(ctx context.Context, log *models.Log) error {

	Logger := logger.CreateFileLoggerWithCtx(ctx)

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

	err := repo.DB.DB().Exec(ctx, query,
		log.Timestamp,
		log.Namespace,
		log.Host,
		log.Service,
		log.Level,
		log.UserID,
		log.SessionID,
		log.TraceID,
		log.Source,
		log.StringNames, log.StringValues,
		log.IntNames, log.IntValues,
		log.FloatNames, log.FloatValues,
		log.BoolNames, log.BoolValues,
	)

	if err != nil {
		Logger.Errorf("Error inserting log: %v", err)
		return err
	}
	return nil
}

func (repo *Repo) SaveBulkLog(ctx context.Context, logs []*models.Log) error {
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	// Bulk insert query must match the single-row insert
	query := `
		INSERT INTO logs (
			_timestamp, _namespace, host, service, level,
			user_id, session_id, trace_id,
			_source,
			string_names, string_values,
			int_names, int_values,
			float_names, float_values,
			bool_names, bool_values
		) VALUES
	`

	// Prepare batch
	batch, err := repo.DB.DB().PrepareBatch(ctx, query)
	if err != nil {
		Logger.Errorf("Error preparing batch: %v", err)
		return err
	}

	for _, log := range logs {
		err = batch.Append(
			log.Timestamp,
			log.Namespace,
			log.Host,
			log.Service,
			log.Level,
			log.UserID,
			log.SessionID,
			log.TraceID,
			log.Source,
			log.StringNames, log.StringValues,
			log.IntNames, log.IntValues,
			log.FloatNames, log.FloatValues,
			log.BoolNames, log.BoolValues,
		)
		if err != nil {
			Logger.Errorf("Error appending to batch: %v", err)
			return err
		}
	}

	// Send batch insert
	if err := batch.Send(); err != nil {
		Logger.Errorf("Error sending batch insert: %v", err)
		return err
	}

	return nil
}
