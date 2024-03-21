package routes

import (
	"context"
	"github.com/soerenchrist/logsync/server/internal/log"
	"net/http"
)

const transactionHeader = "X-Transaction-Id"
const requestIdHeader = "X-Request-Id"

func Scope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		args := make([]any, 0)
		transaction := r.Header.Get(requestIdHeader)
		if transaction != "" {
			args = append(args, "transaction", transaction)
		}

		requestId := r.Header.Get(transactionHeader)
		if requestId != "" {
			args = append(args, "request.id", requestId)
		}

		logger := log.With(args)
		ctx := context.WithValue(r.Context(), "logger", logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
