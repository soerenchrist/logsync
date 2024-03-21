package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
	"github.com/soerenchrist/logsync/server/internal/config"
	"github.com/soerenchrist/logsync/server/internal/files"
	"github.com/soerenchrist/logsync/server/internal/log"
	"github.com/soerenchrist/logsync/server/internal/model"
	"github.com/soerenchrist/logsync/server/internal/routes"
	"net/http"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Error("Could not read config", err)
		return
	}

	logger := log.New(conf.Logging.Level)

	r := chi.NewRouter()
	r.Use(slogchi.New(logger))
	r.Use(middleware.Recoverer)

	db, err := model.CreateDb(conf.Db.Path)
	if err != nil {
		panic(err)
	}

	f := files.New(conf.Files.Path)

	c := routes.NewController(db, r, f)
	c.MapEndpoints()

	log.Info("Server is listening", "url", conf.Url())
	err = http.ListenAndServe(conf.Url(), r)
	if err != nil {
		panic(err)
	}
}
