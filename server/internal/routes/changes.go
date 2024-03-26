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

func (c *Controller) getChanges(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*slog.Logger)
	graphId := chi.URLParam(r, "graphID")

	since := r.URL.Query().Get("since")
	sinceTime, err := parseTime(since)
	if err != nil {
		abort400(w, r, "Could not parse time")
		return
	}
	logger.Debug("Getting changes for graph", "graph", graphId, "since", sinceTime)

	var changes []model.ChangeLogEntry
	tx := c.db.Where("graph_name = ? AND timestamp > ?", graphId, sinceTime).Find(&changes)
	if tx.Error != nil {
		abort500(w, r, err)
		return
	}

	render.JSON(w, r, changes)
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
