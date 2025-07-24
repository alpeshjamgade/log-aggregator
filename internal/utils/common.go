package utils

import (
	"context"
	"fmt"
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

func ExtractClientIDFromMessage(msg string) string {
	re := regexp.MustCompile(`client:\s*([A-Za-z0-9_]+)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func FlattenMap(prefix string, input map[string]any) map[string]any {
	flatMap := make(map[string]any)

	for k, v := range input {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		switch child := v.(type) {
		case map[string]any:
			nested := FlattenMap(fullKey, child)
			for nk, nv := range nested {
				flatMap[nk] = nv
			}
		case map[any]any:
			// convert to map[string]any before flattening
			strMap := make(map[string]any)
			for key, val := range child {
				strMap[fmt.Sprintf("%v", key)] = val
			}
			nested := FlattenMap(fullKey, strMap)
			for nk, nv := range nested {
				flatMap[nk] = nv
			}
		default:
			flatMap[fullKey] = v
		}
	}

	return flatMap
}
