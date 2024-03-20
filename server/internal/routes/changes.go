package routes

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/soerenchrist/logsync/server/internal/model"
	"net/http"
	"strconv"
	"time"
)

func (c *Controller) getChanges(writer http.ResponseWriter, request *http.Request) {
	graphId := chi.URLParam(request, "graphID")

	since := request.URL.Query().Get("since")
	sinceTime, err := parseTime(since)
	if err != nil {
		fmt.Printf("Error: %v", err)
		http.Error(writer, "Could not parse time", http.StatusBadRequest)
		return
	}

	var changes []model.ChangeLogEntry
	tx := c.db.Where("graph_name = ? AND timestamp > ?", graphId, sinceTime).Find(&changes)
	if tx.Error != nil {
		http.Error(writer, "Error in sql query", http.StatusInternalServerError)
		return
	}

	render.JSON(writer, request, changes)
}

func parseTime(since string) (time.Time, error) {
	if since == "" {
		return time.UnixMilli(0), nil
	}

	sinceMillis, err := strconv.ParseInt(since, 10, 64)
	if err != nil {
		return time.UnixMilli(0), err
	}
	return time.UnixMilli(sinceMillis), nil
}
