package app

import (
	"context"
	"go.uber.org/zap"
	"log"
	"log-aggregator/config"
	"log-aggregator/internal/client"
	consumer "log-aggregator/internal/collector"
	"log-aggregator/internal/constants"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/repo"
	"log-aggregator/internal/utils"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Start() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), constants.TraceID, utils.GetUUID())
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	var Logger *zap.SugaredLogger
	if config.LogFile {
		Logger = logger.CreateFileLoggerWithCtx(ctx)
	} else {
		Logger = logger.CreateFileLoggerWithCtx(ctx)
	}

	RmqClient, DB, HttpClient := client.GetClients(ctx)
	Repo := repo.NewRepo(DB, HttpClient)

	Consumer := consumer.NewConsumer(ctx, RmqClient, Repo)

	go func() {
		Logger.Info("starting collector")

		Consumer.Start(ctx)
	}()

	<-ctx.Done()

	Logger.Info("shutting down collector")
}
