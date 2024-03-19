package routes

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/soerenchrist/logsync/server/internal/model"
	"net/http"
	"os"
	"slices"
	"time"
)

func (c *Controller) deleteFile(writer http.ResponseWriter, request *http.Request) {
	graphName := chi.URLParam(request, "graphID")
	fileName := chi.URLParam(request, "fileID")

	err := c.files.Remove(graphName, fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(writer, "Not found", http.StatusNotFound)
		} else {
			http.Error(writer, "Could not delete file", http.StatusInternalServerError)
		}
		return
	}

	timestamp, err := readModifiedDate(request)
	if err != nil {
		http.Error(writer, "Missing modified-date", http.StatusBadRequest)
		return
	}

	fileId := request.FormValue("file-id")
	if fileId == "" {
		http.Error(writer, "Missing file-id", http.StatusBadRequest)
		return
	}

	entry := model.ChangeLogEntry{
		GraphName: graphName,
		FileId:    fileId,
		Operation: model.Deleted,
		Timestamp: timestamp,
		FileName:  fileName,
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

	fileId := request.FormValue("file-id")
	if fileId == "" {
		http.Error(writer, "Missing file-id", http.StatusBadRequest)
		return
	}

	timestamp, err := readModifiedDate(request)
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

	err = c.files.Store(graphName, header.Filename, file)
	if err != nil {
		http.Error(writer, "Could not save file", http.StatusInternalServerError)
		return
	}

	entry := model.ChangeLogEntry{
		GraphName: graphName,
		FileId:    fileId,
		Operation: opType,
		Timestamp: timestamp,
		FileName:  header.Filename,
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

func readModifiedDate(request *http.Request) (time.Time, error) {
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

func validateOperation(operation string) (model.OperationType, error) {
	opType := model.OperationType(operation)
	if slices.Contains(uploadAllowedOperationTypes, opType) {
		return opType, nil
	}

	return "", errors.New(fmt.Sprintf("operation type %s not allowed", operation))
}
