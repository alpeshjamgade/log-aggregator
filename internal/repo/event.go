package repo

import (
	"context"
	"encoding/json"
	"log"
	"log-aggregator/internal/models"
)

func (repo *Repo) SaveEvent(ctx context.Context, eventData *models.Event) error {
	dataJson, err := json.Marshal(eventData.Data)
	if err != nil {
		return err
	}

	err = repo.DB.DB().Exec(ctx,
		"INSERT INTO events(data) VALUES ($1)",
		string(dataJson),
	)

	if err != nil {
		log.Printf("Error inserting event: %v", err)
		return err
	}
	return nil
}
