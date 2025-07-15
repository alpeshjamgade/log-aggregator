package service

import (
	"context"
	"log-aggregator/internal/models"
	"log-aggregator/internal/repo"
)

type IService interface {
	SaveEvent(ctx context.Context, event *models.Event) error
}

type Service struct {
	repo repo.IRepo
}

func NewService(repo repo.IRepo) IService {
	return &Service{repo: repo}
}
