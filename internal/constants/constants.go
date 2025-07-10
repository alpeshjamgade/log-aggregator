package constants

import (
	"github.com/gorilla/sessions"
	"log-aggregator/config"
)

const (
	TraceID     = "trace_id"
	Service     = "service"
	ServiceName = "log-aggregator"
	Empty       = ""
)

var (
	CookieStore = sessions.NewCookieStore([]byte(config.SessionKey))
)
