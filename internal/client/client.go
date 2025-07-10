package client

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.elastic.co/apm/module/apmhttp/v2"
	"log-aggregator/config"
	"log-aggregator/internal/client/db"
	"log-aggregator/internal/logger"
	"net/http"
	"os"
	"time"
)

var (
	HttpClient *http.Client     = nil
	RMQClient  *amqp.Connection = nil
	DbClient   db.DB            = nil
)

func GetClients(ctx context.Context) (*amqp.Connection, db.DB, *http.Client) {

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

	dbUrl := config.DatabaseURL
	DbClient = db.NewPostgresDB(dbUrl)
	err = DbClient.Connect(ctx)
	if err != nil {
		Logger.Panic(err)
		os.Exit(1)
	}

	RMQClient, err = amqp.Dial(config.RMQUrl)
	if err != nil {
		Logger.DPanicf("Failed to connect to RabbitMQ: %v", err)
	}

	Logger.Info("Connected to RabbitMQ")

	return RMQClient, DbClient, HttpClient
}
