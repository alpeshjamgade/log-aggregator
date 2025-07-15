package db

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type DB interface {
	Connect(ctx context.Context) error
	Disconnect() error
	DB() driver.Conn
}
