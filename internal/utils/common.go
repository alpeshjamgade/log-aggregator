package utils

import (
	"context"
	"github.com/google/uuid"
)

func GetUUID() string {
	return uuid.New().String()
}

func ContextWithValueIfNotPresent(ctx context.Context, key string, value string) context.Context {
	if ctx.Value(key) == nil {
		ctx = context.WithValue(ctx, key, value)
	}

	return ctx
}
