package routes

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Controller struct {
	db     *gorm.DB
	router *chi.Mux
}

func NewController(db *gorm.DB, router *chi.Mux) *Controller {
	c := &Controller{
		db:     db,
		router: router,
	}
	return c
}

func (c *Controller) MapEndpoints() {
	c.router.Get("/changes/{graphID}", c.getChanges)
	c.router.Post("/upload", c.uploadFile)
}
