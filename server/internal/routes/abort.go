package routes

import (
	"log/slog"
	"net/http"
)

func abort500(w http.ResponseWriter, r *http.Request, err error) {
	logger := r.Context().Value("logger").(*slog.Logger)
	logger.Error("Could not query database", "error", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
