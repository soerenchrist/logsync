package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/soerenchrist/logsync/server/internal/files"
	"gorm.io/gorm"
)

type Controller struct {
	db     *gorm.DB
	router *chi.Mux
	files  files.FileStore
}

func NewController(db *gorm.DB, router *chi.Mux, f files.FileStore) *Controller {
	c := &Controller{
		db:     db,
		router: router,
		files:  f,
	}
	return c
}

func (c *Controller) MapEndpoints() {
	c.router.Get("/{graphID}/changes", c.getChanges)
	c.router.Post("/{graphID}/upload", c.uploadFile)
	c.router.Post("/{graphID}/delete/{fileID}", c.deleteFile)
	c.router.Post("/{graphID}/content/{fileID}", c.content)
}
