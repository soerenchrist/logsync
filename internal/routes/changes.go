package routes

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type ChangesController struct {
	db     *gorm.DB
	router *chi.Mux
}

func NewChangesController(db *gorm.DB, router *chi.Mux) *ChangesController {
	c := &ChangesController{
		db:     db,
		router: router,
	}
	return c
}

func (c *ChangesController) MapEndpoints() {
	c.router.Get("/changes", c.getChanges)
}

func (c *ChangesController) getChanges(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Changes called")
}
