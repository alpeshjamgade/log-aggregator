package repo

import (
	"context"
	"log-aggregator/internal/client/db"
	"log-aggregator/internal/models"
	"net/http"
)

type IRepo interface {
	SaveEvent(ctx context.Context, event *models.Event) error
}

type Repo struct {
	DB         db.DB
	HttpClient *http.Client
}

func NewRepo(db db.DB, httpClient *http.Client) *Repo {
	return &Repo{DB: db, HttpClient: httpClient}
}
