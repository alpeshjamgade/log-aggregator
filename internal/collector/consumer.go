package consumer

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"log"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/repo"
	"sync"
	"time"
)

type IConsumer interface {
}

type Consumer struct {
	rmqConn    *amqp.Connection
	rmqQueue   *amqp.Queue
	rmqChannel *amqp.Channel
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup

	// Metrics
	processedCount int64
	errorCount     int64
	lastBatchTime  time.Time

	// Channels
	messageChan chan LogEntry
	batchChan   chan []LogEntry

	// Database
	repo *repo.Repo

	logger *zap.SugaredLogger
}

// LogEntry represents a parsed log entry
type LogEntry struct {
	Timestamp  time.Time         `json:"timestamp"`
	Level      string            `json:"level"`
	Message    string            `json:"message"`
	Source     string            `json:"source"`
	Host       string            `json:"host"`
	Fields     map[string]string `json:"fields"`
	RawMessage string            `json:"raw_message"`
}

func NewConsumer(ctx context.Context, rmqClient *amqp.Connection, repo *repo.Repo) *Consumer {
	ctx, cancel := context.WithCancel(context.WithValue(ctx, "rmqClient", rmqClient))
	defer cancel()
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	channel, err := rmqClient.Channel()
	if err != nil {
		cancel()
		Logger.Fatal("Failed to open channel", zap.Error(err))
	}

	return &Consumer{ctx: ctx, rmqConn: rmqClient, repo: repo, rmqChannel: channel}
}

func (consumer *Consumer) Start(ctx context.Context) {
	consumer.wg.Add(1)
	go consumer.startWorker()
}

func (consumer *Consumer) startWorker() {
	defer consumer.wg.Done()

	var receiverChannel chan struct{}

	msgs, err := consumer.rmqChannel.Consume(
		"logs", // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		consumer.logger.Fatal("Failed to register a consumer", zap.Error(err))
	}

	go func() {
		for d := range msgs {
			consumer.repo.InsertLog(string(d.Body))

			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-receiverChannel
}
