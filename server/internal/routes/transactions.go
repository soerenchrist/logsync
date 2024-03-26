package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/soerenchrist/logsync/server/internal/model"
	"log/slog"
	"net/http"
	"strconv"
)

type pageOptions struct {
	page int
	size int
}

func (p pageOptions) skip() int {
	page := p.page - 1
	return page * p.size
}

func (c *Controller) getChangesInTransaction(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*slog.Logger)
	transactionId := chi.URLParam(r, "transactionID")
	page := readPageOptions(r)
	logger.Debug("Getting changes for transaction", "id", transactionId)

	var changes []model.ChangeLogEntry
	tx := c.db.Where("transaction_id = ?", transactionId).
		Order("timestamp desc").
		Limit(page.size).
		Offset(page.skip()).
		Find(&changes)
	if tx.Error != nil {
		abort500(w, r, tx.Error)
		return
	}

	logger.Debug("Found changes for transaction", "count", len(changes))

	render.JSON(w, r, changes)
}

func (c *Controller) getTransactions(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*slog.Logger)
	page := readPageOptions(r)
	logger.Debug("Getting transactions", "page", page.page, "size", page.size)

	var result []struct {
		From  string `json:"from"`
		To    string `json:"to"`
		Count int    `json:"count"`
		Id    string `json:"id"`
	}
	tx := c.db.Raw(`
		SELECT transactions.transaction_id as id,
		(SELECT max(timestamp) FROM change_log_entries WHERE transaction_id = transactions.transaction_id) as 'to',
		(SELECT min(timestamp) FROM change_log_entries WHERE transaction_id = transactions.transaction_id) as 'from',
        (SELECT count(*) FROM change_log_entries WHERE transaction_id = transactions.transaction_id)       as [count]
	FROM (SELECT DISTINCT transaction_id FROM change_log_entries) as transactions
		ORDER BY [to] DESC 
		LIMIT ? OFFSET ?
	`, page.size, page.skip()).Scan(&result)
	if tx.Error != nil {
		abort500(w, r, tx.Error)
		return
	}

	logger.Debug("Found transactions", "count", len(result))

	render.JSON(w, r, result)
}

func readPageOptions(r *http.Request) pageOptions {
	size := 10
	page := 1
	var err error
	sizeQuery := r.URL.Query().Get("size")
	if sizeQuery != "" {
		size, err = strconv.Atoi(sizeQuery)
		if err != nil || size < 1 {
			size = 10
		}
	}

	err = nil
	pageQuery := r.URL.Query().Get("page")
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil || page <= 1 {
			page = 1
		}
	}

	return pageOptions{page: page, size: size}
}
