package routes

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/soerenchrist/logsync/server/internal/model"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"time"
)

func (c *Controller) content(w http.ResponseWriter, r *http.Request) {
	graphName := chi.URLParam(r, "graphID")
	fileId := chi.URLParam(r, "fileID")

	mapping, err := c.getMapping(fileId)
	if err != nil {
		abort500(w, r, err)
		return
	}

	data, err := c.files.Content(graphName, mapping.FileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			abort404(w, r)
		} else {
			abort500(w, r, err)
		}
		return
	}

	render.Data(w, r, data)
}

func (c *Controller) getMapping(fileId string) (model.FileMapping, error) {
	fileId, err := url.PathUnescape(fileId)
	if err != nil {
		return model.FileMapping{}, err
	}

	var fileMapping model.FileMapping
	tx := c.db.Where("file_id = ?", fileId).First(&fileMapping)
	if tx.Error != nil {
		return model.FileMapping{}, tx.Error
	}

	return fileMapping, nil
}

func (c *Controller) getOrCreateMapping(fileId string) (model.FileMapping, error) {
	fileId, err := url.PathUnescape(fileId)
	if err != nil {
		return model.FileMapping{}, err
	}

	var found model.FileMapping
	tx := c.db.Where("file_id = ?", fileId).First(&found)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return model.FileMapping{}, tx.Error
		}
	} else {
		return found, nil
	}
	mapping := model.FileMapping{
		FileId:   fileId,
		FileName: uuid.New().String(),
	}
	tx = c.db.Create(&mapping)
	return mapping, tx.Error
}

func (c *Controller) removeMapping(fileId string) error {
	fileId, err := url.PathUnescape(fileId)
	if err != nil {
		return err
	}
	tx := c.db.Delete(&model.FileMapping{
		FileId: fileId,
	})
	return tx.Error
}

func (c *Controller) deleteFile(w http.ResponseWriter, r *http.Request) {
	graphName := chi.URLParam(r, "graphID")
	fileName := chi.URLParam(r, "fileID")

	transaction := r.URL.Query().Get("ta_id")
	if transaction == "" {
		abort400(w, r, "ta_id query param is missing")
		return
	}

	timestamp, err := readModifiedDateFromQuery(r)
	if err != nil {
		abort400(w, r, "modified_date query param is missing")
		return
	}

	// check of duplicate entry
	var existing []model.ChangeLogEntry
	c.db.Where("timestamp = ? AND file_id = ? and graph_name = ?", timestamp, fileName, graphName).Find(&existing)
	if len(existing) > 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	mapping, err := c.getOrCreateMapping(fileName)
	if err != nil {
		abort500(w, r, err)
		return
	}

	err = c.files.Remove(graphName, mapping.FileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			abort404(w, r)
		} else {
			abort500(w, r, err)
		}
		return
	}

	err = c.removeMapping(fileName)
	if err != nil {
		abort500(w, r, err)
		return
	}
	entry := model.ChangeLogEntry{
		GraphName:     graphName,
		FileId:        fileName,
		Operation:     model.Deleted,
		Timestamp:     timestamp,
		TransactionId: transaction,
	}
	tx := c.db.Create(entry)
	if tx.Error != nil {
		abort500(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) uploadFile(w http.ResponseWriter, r *http.Request) {
	graphName := chi.URLParam(r, "graphID")
	err := r.ParseMultipartForm(10 << 20) // max of 10MB
	if err != nil {
		abort400(w, r, "Expected multipart body")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		abort400(w, r, "Expected file parameter in form")
		return
	}
	defer file.Close()

	transaction := r.FormValue("ta-id")
	if transaction == "" {
		abort400(w, r, "Expected ta-id parameter in form")
		return
	}

	timestamp, err := readModifiedDateFromForm(r)
	if err != nil {
		abort400(w, r, "Could not parse modified-date")
		return
	}

	operation := r.FormValue("operation")
	if operation == "" {
		abort400(w, r, "Missing operation parameter in form")
		return
	}
	opType, err := validateOperation(operation)
	if err != nil {
		abort400(w, r, fmt.Sprintf("Invalid operation type, allowed values: %v", uploadAllowedOperationTypes))
		return
	}

	// check of duplicate entry
	var existing []model.ChangeLogEntry
	c.db.Where("timestamp = ? AND file_id = ? and graph_name = ?", timestamp, header.Filename, graphName).Find(&existing)
	if len(existing) > 0 {
		w.WriteHeader(http.StatusCreated)
		return
	}

	mapping, err := c.getOrCreateMapping(header.Filename)
	if err != nil {
		abort500(w, r, err)
		return
	}

	err = c.files.Store(graphName, mapping.FileName, file)
	if err != nil {
		abort500(w, r, err)
		return
	}

	entry := model.ChangeLogEntry{
		GraphName:     graphName,
		FileId:        header.Filename,
		Operation:     opType,
		Timestamp:     timestamp,
		TransactionId: transaction,
	}
	tx := c.db.Create(entry)
	if tx.Error != nil {
		abort500(w, r, tx.Error)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

var uploadAllowedOperationTypes = []model.OperationType{
	model.Modified, model.Created,
}

func readModifiedDateFromForm(request *http.Request) (time.Time, error) {
	modifiedDate := request.FormValue("modified-date")
	if modifiedDate == "" {
		return time.Now(), nil
	}
	timestamp, err := time.Parse(time.RFC3339, modifiedDate)
	if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

func readModifiedDateFromQuery(request *http.Request) (time.Time, error) {
	modifiedDate := request.URL.Query().Get("modified_date")
	if modifiedDate == "" {
		return time.Now(), nil
	}

	millis, err := strconv.ParseInt(modifiedDate, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	timestamp := time.UnixMilli(millis)
	return timestamp, nil
}

func validateOperation(operation string) (model.OperationType, error) {
	opType := model.OperationType(operation)
	if slices.Contains(uploadAllowedOperationTypes, opType) {
		return opType, nil
	}

	return "", errors.New(fmt.Sprintf("operation type %s not allowed", operation))
}
