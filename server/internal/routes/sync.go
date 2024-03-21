package routes

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/soerenchrist/logsync/server/internal/model"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"
)

func (c *Controller) content(writer http.ResponseWriter, request *http.Request) {
	graphName := chi.URLParam(request, "graphID")
	fileName := chi.URLParam(request, "fileID")

	data, err := c.files.Content(graphName, fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("File %s not found", fileName)
			http.Error(writer, "Not found", http.StatusNotFound)
		} else {
			fmt.Printf("Could not read file %s, %v", fileName, err)
			http.Error(writer, "Could not read file", http.StatusInternalServerError)
		}
		return
	}

	render.Data(writer, request, data)
}

func (c *Controller) deleteFile(writer http.ResponseWriter, request *http.Request) {
	graphName := chi.URLParam(request, "graphID")
	fileName := chi.URLParam(request, "fileID")

	transaction := request.URL.Query().Get("ta_id")
	if transaction == "" {
		http.Error(writer, "ta_id query is missing", http.StatusBadRequest)
		return
	}

	timestamp, err := readModifiedDateFromQuery(request)
	if err != nil {
		http.Error(writer, "Missing modified_date query param", http.StatusBadRequest)
		return
	}

	// check of duplicate entry
	var existing []model.ChangeLogEntry
	c.db.Where("timestamp = ? AND file_id = ? and graph_name = ?", timestamp, fileName, graphName).Find(&existing)
	if len(existing) > 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	err = c.files.Remove(graphName, fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(writer, "Not found", http.StatusNotFound)
		} else {
			http.Error(writer, "Could not delete file", http.StatusInternalServerError)
		}
		return
	}

	entry := model.ChangeLogEntry{
		GraphName:   graphName,
		FileId:      fileName,
		Operation:   model.Deleted,
		Timestamp:   timestamp,
		FileName:    fileName,
		Transaction: transaction,
	}
	tx := c.db.Create(entry)
	if tx.Error != nil {
		http.Error(writer, "Failed to save entry", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (c *Controller) uploadFile(writer http.ResponseWriter, request *http.Request) {
	graphName := chi.URLParam(request, "graphID")
	err := request.ParseMultipartForm(10 << 20) // max of 10MB
	if err != nil {
		http.Error(writer, "Expected multipart body", http.StatusBadRequest)
		return
	}

	file, header, err := request.FormFile("file")
	if err != nil {
		http.Error(writer, "Could not read request", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	transaction := request.FormValue("ta-id")
	if transaction == "" {
		http.Error(writer, "Missing ta-id", http.StatusBadRequest)
		return
	}

	timestamp, err := readModifiedDateFromForm(request)
	if err != nil {
		http.Error(writer, "Could not parse modified-date", http.StatusBadRequest)
		return
	}

	operation := request.FormValue("operation")
	if operation == "" {
		http.Error(writer, "Missing operation", http.StatusBadRequest)
		return
	}
	opType, err := validateOperation(operation)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid operation type, allowed values: %v", uploadAllowedOperationTypes), http.StatusBadRequest)
		return
	}

	// check of duplicate entry
	var existing []model.ChangeLogEntry
	c.db.Where("timestamp = ? AND file_id = ? and graph_name = ?", timestamp, header.Filename, graphName).Find(&existing)
	if len(existing) > 0 {
		writer.WriteHeader(http.StatusCreated)
		return
	}

	err = c.files.Store(graphName, header.Filename, file)
	if err != nil {
		http.Error(writer, "Could not save file", http.StatusInternalServerError)
		return
	}

	entry := model.ChangeLogEntry{
		GraphName:   graphName,
		FileId:      header.Filename,
		Operation:   opType,
		Timestamp:   timestamp,
		FileName:    header.Filename,
		Transaction: transaction,
	}
	tx := c.db.Create(entry)
	if tx.Error != nil {
		http.Error(writer, "Could not save file", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
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
