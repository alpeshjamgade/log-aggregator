package db

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"log-aggregator/internal/logger"
	"time"
)

type ClickHouseDB struct {
	Conn        driver.Conn
	DatabaseURL string
}

func NewClickHouseDB(databaseURL string) *ClickHouseDB {
	return &ClickHouseDB{DatabaseURL: databaseURL}
}

func (c *ClickHouseDB) DB() driver.Conn { return c.Conn }

func (c *ClickHouseDB) Connect(ctx context.Context) error {
	var err error
	var count int8

	Logger := logger.CreateLoggerWithCtx(ctx)

	dbUrl := c.DatabaseURL
	for {
		c.Conn, _ = clickhouse.Open(&clickhouse.Options{
			Addr: []string{c.DatabaseURL},
			Auth: clickhouse.Auth{
				Database: "log_aggregator",
			},
			DialTimeout:      time.Duration(10) * time.Second,
			MaxOpenConns:     100,
			MaxIdleConns:     10,
			ConnMaxLifetime:  time.Duration(2) * time.Minute,
			ConnOpenStrategy: clickhouse.ConnOpenInOrder,
			BlockBufferSize:  10,
		})
		err = c.Conn.Ping(ctx)

		if err != nil {
			Logger.Errorf("Error connecting to clickhouse: %v", err)
			count++
		} else {
			Logger.Infof("connected to clickhouse at %s", dbUrl)
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

func (c *ClickHouseDB) Disconnect() error {
	return c.Conn.Close()
}
