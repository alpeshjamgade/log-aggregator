package repo

import (
	"log"
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

func (repo *Repo) InsertEvent(event string) error {
	_, err := repo.DB.DB().Exec(
		"INSERT INTO events(data) VALUES ($1)",
		event,
	)

	if err != nil {
		log.Printf("Error inserting event: %v", err)
		return err
	}
	return nil
}
