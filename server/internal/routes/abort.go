package routes

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type apiError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func abort500(w http.ResponseWriter, r *http.Request, err error) {
	logger := r.Context().Value("logger").(*slog.Logger)
	logger.Error("An error occurred", "error", err)
	abort(w, r, 500, "An error occurred")
}

func abort400(w http.ResponseWriter, r *http.Request, message string) {
	abort(w, r, 400, message)
}

func abort404(w http.ResponseWriter, r *http.Request) {
	abort(w, r, 404, "Not found")
}

func abort(w http.ResponseWriter, r *http.Request, status int, error string) {
	render.Status(r, status)
	render.JSON(w, r, apiError{
		Code:  status,
		Error: error,
	})
}
