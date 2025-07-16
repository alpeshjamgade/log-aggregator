package utils

import (
	"context"
	"github.com/google/uuid"
	"regexp"
	"strings"
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

func ExtractJSONFromLog(s string) string {
	// Optional: strip ANSI escape codes
	reAnsi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	s = reAnsi.ReplaceAllString(s, "")

	// Find first `{` and extract everything from there
	start := strings.Index(s, "{")
	if start == -1 {
		return s // fallback
	}
	return s[start:]
}
