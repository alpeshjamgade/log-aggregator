package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq"
	"log-aggregator/internal/logger"
	"time"
)

type PostgresDB struct {
	Sql         *sql.DB
	DatabaseURL string
}

func NewPostgresDB(url string) *PostgresDB {
	return &PostgresDB{
		DatabaseURL: url,
	}
}

func (p *PostgresDB) DB() *sql.DB { return p.Sql }
func (p *PostgresDB) Connect(ctx context.Context) error {
	var err error
	var count int8
	Logger := logger.CreateLoggerWithCtx(ctx)

	dbUrl := p.DatabaseURL
	for {
		//p.Sql, err = sqlx.Connect("postgres", dbUrl)
		p.Sql, _ = apmsql.Open("postgres", dbUrl)
		err = p.Sql.Ping()

		if err != nil {
			Logger.Errorf("Error connecting to Postgres: %v", err)
			count++
		} else {
			Logger.Infof("connected to postgres at %s", dbUrl)
			p.Sql.SetMaxOpenConns(5000)
			p.Sql.SetMaxIdleConns(1000)
			p.Sql.SetConnMaxLifetime(2 * time.Minute)
			break
		}

		if count > 5 {
			Logger.Errorf(err.Error())
			return err
		}
		Logger.Warnf("Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)

	}

	return nil
}

func (p *PostgresDB) Disconnect() error { return p.Sql.Close() }
