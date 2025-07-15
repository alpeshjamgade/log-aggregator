package client

import (
	"context"
	"go.elastic.co/apm/module/apmhttp/v2"
	"log-aggregator/config"
	"log-aggregator/internal/client/db"
	"log-aggregator/internal/logger"
	"net/http"
	"os"
	"time"
)

var (
	HttpClient *http.Client = nil
	DbClient   db.DB
)

func GetClients(ctx context.Context) (db.DB, *http.Client) {

	Logger := logger.CreateFileLoggerWithCtx(ctx)

	var err error

	HttpClient = apmhttp.WrapClient(&http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        500,
			MaxConnsPerHost:     500,
			MaxIdleConnsPerHost: 500,
			IdleConnTimeout:     20 * time.Second,
		},
	})

	dbUrl := config.ClickHouseDBAddr
	DbClient = db.NewClickHouseDB(dbUrl)
	err = DbClient.Connect(ctx)
	if err != nil {
		Logger.Panic(err)
		os.Exit(1)
	}

	return DbClient, HttpClient
}
