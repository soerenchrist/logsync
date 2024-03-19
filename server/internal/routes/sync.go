package routes

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/soerenchrist/logsync/server/internal/model"
	"net/http"
	"os"
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

	writer.WriteHeader(http.StatusNoContent)
}

func (c *Controller) uploadFile(writer http.ResponseWriter, request *http.Request) {
	graphName := chi.URLParam(request, "graphID")
	err := request.ParseMultipartForm(10 << 20) // max of 10MB
	if err != nil {
		http.Error(writer, "Could not read request", http.StatusInternalServerError)
		return
	}

	file, header, err := request.FormFile("file")
	if err != nil {
		http.Error(writer, "Could not read request", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var timestamp time.Time
	modifiedDate := request.FormValue("modified-date")
	if modifiedDate == "" {
		timestamp = time.Now()
	} else {
		timestamp, err = time.Parse(time.RFC3339, modifiedDate)
		if err != nil {
			http.Error(writer, "Could not parse modified-date", http.StatusBadRequest)
			return
		}
	}

	fileId := request.FormValue("file-id")
	if fileId == "" {
		http.Error(writer, "Missing file-id", http.StatusBadRequest)
		return
	}

	operation := request.FormValue("operation")
	if operation == "" {
		http.Error(writer, "Missing operation", http.StatusBadRequest)
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
		Operation: model.OperationType(operation),
		Timestamp: timestamp,
		FileName:  header.Filename,
	}
	tx := c.db.Create(entry)
	if tx.Error != nil {
		http.Error(writer, "Could not save file", http.StatusInternalServerError)
		return
	}

	fmt.Printf("File: %s\n", header.Filename)
	fmt.Printf("Size: %d\n", header.Size)
	fmt.Printf("Header: %v\n", header.Header)
	fmt.Printf("Modified: %s\n", modifiedDate)

	writer.WriteHeader(http.StatusCreated)
}
