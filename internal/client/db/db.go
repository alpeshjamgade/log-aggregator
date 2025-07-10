package db

import (
	"context"
	"database/sql"
)

type DB interface {
	Connect(ctx context.Context) error
	Disconnect() error
	DB() *sql.DB
}
