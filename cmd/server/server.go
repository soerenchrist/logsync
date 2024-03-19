package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/soerenchrist/logsync/internal/model"
	"github.com/soerenchrist/logsync/internal/routes"
	"net/http"
)

func Start() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := model.CreateDb("test.db")
	if err != nil {
		return err
	}

	c := routes.NewChangesController(db, r)
	c.MapEndpoints()

	return http.ListenAndServe(":3000", r)
}
