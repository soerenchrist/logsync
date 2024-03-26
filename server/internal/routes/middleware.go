package routes

import (
	"context"
	"github.com/soerenchrist/logsync/server/internal/config"
	"github.com/soerenchrist/logsync/server/internal/log"
	"net/http"
)

const transactionHeader = "X-Transaction-Id"
const requestIdHeader = "X-Request-Id"
const apiTokenHeader = "X-Api-Token"

func Scope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		args := make([]any, 0)
		transaction := r.Header.Get(transactionHeader)
		if transaction != "" {
			args = append(args, "transaction", transaction)
		}

		requestId := r.Header.Get(requestIdHeader)
		if requestId != "" {
			args = append(args, "request.id", requestId)
		}

		logger := log.With(args)
		ctx := context.WithValue(r.Context(), "logger", logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateApiTokenMiddleware(conf config.Config) func(handler http.Handler) http.Handler {
	if conf.Server.ApiToken == "" {
		return noop
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(apiTokenHeader)
			if token != conf.Server.ApiToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func noop(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
