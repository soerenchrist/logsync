package routes

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/soerenchrist/logsync/server/internal/model"
	"net/http"
	"time"
)

func (c *Controller) getChanges(writer http.ResponseWriter, request *http.Request) {
	graphId := chi.URLParam(request, "graphID")

	since := request.URL.Query().Get("since")
	sinceTime, err := parseTime(since)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	var changes []model.ChangeLogEntry
	tx := c.db.Where("graph_name = ? AND timestamp > ?", graphId, sinceTime).Find(&changes)
	if tx.Error != nil {
		fmt.Printf("Error while querying: %v", err)
		return
	}

	for i, change := range changes {
		fmt.Printf("Change %d: %v", i, change)
	}
}

func parseTime(since string) (time.Time, error) {
	if since == "" {
		return time.UnixMilli(0), nil
	}

	return time.Parse(time.RFC3339, since)
}
