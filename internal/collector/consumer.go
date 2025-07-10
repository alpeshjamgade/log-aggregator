package consumer

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
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
	return &Consumer{ctx: ctx, rmqConn: rmqClient, repo: repo}
}

func (consumer *Consumer) Start(ctx context.Context) {

}
