package service

import (
	"context"
	"log-aggregator/internal/models"
)

func (svc *Service) SaveEvent(ctx context.Context, event *models.Event) error {
	err := svc.repo.SaveEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
