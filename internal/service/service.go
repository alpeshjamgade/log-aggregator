package service

import (
	"context"
	"log-aggregator/internal/models"
	"log-aggregator/internal/repo"
)

type IService interface {
	SaveLog(ctx context.Context, log *models.FluentBitReq) error
	SaveBulkLog(ctx context.Context, log []*models.FluentBitReq) error
	SaveBulkLogV2(ctx context.Context, log []map[string]any) error
}

type Service struct {
	repo repo.IRepo
}

func NewService(repo repo.IRepo) IService {
	return &Service{repo: repo}
}
