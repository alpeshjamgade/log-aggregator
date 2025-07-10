package repo

import (
	"log-aggregator/internal/client/db"
	"net/http"
)

type Repo struct {
	DB         db.DB
	HttpClient *http.Client
}

func NewRepo(db db.DB, httpClient *http.Client) *Repo {
	return &Repo{DB: db, HttpClient: httpClient}
}
