package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/soerenchrist/logsync/server/internal/model"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (c *Controller) getChanges(writer http.ResponseWriter, request *http.Request) {
	logger := request.Context().Value("logger").(*slog.Logger)
	graphId := chi.URLParam(request, "graphID")

	since := request.URL.Query().Get("since")
	sinceTime, err := parseTime(since)
	if err != nil {
		logger.Error("Could not parse time", "error", err)
		http.Error(writer, "Could not parse time", http.StatusBadRequest)
		return
	}
	logger.Debug("Getting changes for graph", "graph", graphId, "since", sinceTime)

	var changes []model.ChangeLogEntry
	tx := c.db.Where("graph_name = ? AND timestamp > ?", graphId, sinceTime).Find(&changes)
	if tx.Error != nil {
		logger.Error("Could not query database", "error", err)
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
