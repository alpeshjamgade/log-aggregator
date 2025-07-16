package repo

import (
	"context"
	"log-aggregator/internal/client/db"
	"log-aggregator/internal/models"
	"net/http"
)

type IRepo interface {
	SaveLog(ctx context.Context, log *models.Log) error
	SaveBulkLog(ctx context.Context, log []*models.Log) error
}

type Repo struct {
	DB         db.DB
	HttpClient *http.Client
}

func NewRepo(db db.DB, httpClient *http.Client) *Repo {
	return &Repo{DB: db, HttpClient: httpClient}
}
