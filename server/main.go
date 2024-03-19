package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/soerenchrist/logsync/server/internal/model"
	"github.com/soerenchrist/logsync/server/internal/routes"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := model.CreateDb("test.db")
	if err != nil {
		panic(err)
	}

	c := routes.NewController(db, r)
	c.MapEndpoints()

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
