package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log-aggregator/internal/constants"
	"log-aggregator/internal/utils"
	"net/http"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (cw *CustomResponseWriter) WriteHeader(status int) {
	if cw.status == 0 {
		cw.status = status
	}

	cw.ResponseWriter.WriteHeader(status)
}

func PrettyPrint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := utils.ContextWithValueIfNotPresent(r.Context(), constants.TraceID, utils.GetUUID())

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return
		}

		// Reassign r.Body so it can still be read later by ReadJSON()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Pretty-print the raw JSON
		var prettyBody bytes.Buffer
		if err := json.Indent(&prettyBody, bodyBytes, "", "  "); err != nil {
			fmt.Println("Raw body (not JSON):\n", string(bodyBytes))
		} else {
			fmt.Println("Pretty JSON Body:\n", prettyBody.String())
		}

		cw := &CustomResponseWriter{ResponseWriter: w}
		next.ServeHTTP(cw, r.WithContext(ctx))

	})
}
