package remote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/soerenchrist/logsync/client/internal/config"
	"github.com/soerenchrist/logsync/client/internal/graph"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Remote struct {
	conf config.Config
}

type ChangeLogEntry struct {
	GraphName   string    `json:"graph_name"`
	FileId      string    `json:"file_id"`
	Timestamp   time.Time `json:"timestamp"`
	Transaction string    `json:"transaction"`
	Operation   string    `json:"operation"`
}

func New(conf config.Config) *Remote {
	return &Remote{conf: conf}
}

func (r *Remote) GetChanges(graphName string, since time.Time) ([]ChangeLogEntry, error) {
	url := fmt.Sprintf("%s/%s/changes?since=%d", r.conf.Server.Host, graphName, since.UnixMilli())
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("status code is no success")
	}

	var entries []ChangeLogEntry
	err = json.Unmarshal(body, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *Remote) GetContent(graphName string, fileId string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/content/%s", r.conf.Server.Host, graphName, fileId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("no success status code: %d", resp.StatusCode))
	}

	return io.ReadAll(resp.Body)
}

func (r *Remote) DeleteFile(graphName string, file graph.File, transaction string) error {
	url := fmt.Sprintf("%s/%s/delete/%s?ta_id=%s&modified_date=%d", r.conf.Server.Host, graphName, file.Id, transaction, file.LastChange.UnixMilli())
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("no success status code: %d", resp.StatusCode))
	}

	return nil
}

func (r *Remote) UploadFile(graphName string, file graph.File, transaction, operation string) error {
	url := fmt.Sprintf("%s/%s/upload", r.conf.Server.Host, graphName)
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	fileWriter, err := mw.CreateFormFile("file", file.Id)
	if err != nil {
		return err
	}
	contents, err := os.ReadFile(file.Path)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(contents)
	if err != nil {
		return err
	}

	err = addFormField(mw, "ta-id", transaction)
	if err != nil {
		return err
	}

	err = addFormField(mw, "operation", operation)
	if err != nil {
		return err
	}

	err = addFormField(mw, "modified-date", file.LastChange.Format(time.RFC3339))
	if err != nil {
		return err
	}
	err = mw.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Message: %s\n", body)
		fmt.Printf("Status-code: %d\n", resp.StatusCode)
		return errors.New(fmt.Sprintf("no success status code: %d", resp.StatusCode))
	}
	return nil
}

func addFormField(mw *multipart.Writer, fieldName, content string) error {
	writer, err := mw.CreateFormField(fieldName)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(content))
	if err != nil {
		return err
	}
	return nil
}
