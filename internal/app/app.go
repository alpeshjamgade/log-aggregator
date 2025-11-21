package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"log-aggregator/config"
	"log-aggregator/internal/client"
	"log-aggregator/internal/constants"
	"log-aggregator/internal/handler"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/repo"
	"log-aggregator/internal/service"
	"log-aggregator/internal/utils"
	"net/http"
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
		Logger = logger.CreateLoggerWithCtx(ctx)
	}

	Router := GetRouter()
	DB, HttpClient := client.GetClients(ctx)

	Repo := repo.NewRepo(DB, HttpClient)
	Service := service.NewService(Repo)
	Handler := handler.NewHandler(Service)
	Handler.SetupRoutes(Router)

	go func() {
		Logger.Infof("starting server on http://0.0.0.0:%s", config.HttpPort)
		http.ListenAndServe(fmt.Sprintf(":%s", config.HttpPort), Router)
	}()

	<-ctx.Done()

	Logger.Info("shutting down server")
}
