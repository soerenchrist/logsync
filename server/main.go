package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/soerenchrist/logsync/server/internal/config"
	"github.com/soerenchrist/logsync/server/internal/files"
	"github.com/soerenchrist/logsync/server/internal/model"
	"github.com/soerenchrist/logsync/server/internal/routes"
	"net/http"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Could not read config: %v", err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := model.CreateDb(conf.Db.Path)
	if err != nil {
		panic(err)
	}

	f := files.New(conf.Files.Path)

	c := routes.NewController(db, r, f)
	c.MapEndpoints()

	err = http.ListenAndServe(conf.Url(), r)
	if err != nil {
		panic(err)
	}
}
