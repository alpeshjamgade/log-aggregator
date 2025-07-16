package service

import (
	"context"
	"log-aggregator/internal/models"
	"log-aggregator/internal/repo"
)

type IService interface {
	SaveLog(ctx context.Context, log *models.RawLog) error
}

type Service struct {
	repo repo.IRepo
}

func NewService(repo repo.IRepo) IService {
	return &Service{repo: repo}
}
