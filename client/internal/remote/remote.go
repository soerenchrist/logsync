package remote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/soerenchrist/logsync/client/internal/config"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type ChangesRequest struct {
	conf config.Config
}

type ContentRequest struct {
	conf config.Config
}

type ChangeLogEntry struct {
	GraphName   string    `json:"graph_name"`
	FileId      string    `json:"file_id"`
	Timestamp   time.Time `json:"timestamp"`
	Transaction string    `json:"transaction"`
	Operation   string    `json:"operation"`
}

func NewChangesRequest(conf config.Config) ChangesRequest {
	return ChangesRequest{conf: conf}
}

func NewContentRequest(conf config.Config) ContentRequest {
	return ContentRequest{conf: conf}
}

func (r ChangesRequest) Send(graphName string, since time.Time) ([]ChangeLogEntry, error) {
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

func (r ContentRequest) Send(graphName string, fileId string) ([]byte, error) {
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

type request struct {
	config      config.Config
	graphName   string
	transaction string
	operation   string
}

type UploadRequest struct {
	request
}

type DeleteRequest struct {
	request
}

func NewUploadRequest(conf config.Config, graphName, transaction, operation string) UploadRequest {
	return UploadRequest{
		request: request{
			config:      conf,
			graphName:   graphName,
			transaction: transaction,
			operation:   operation,
		},
	}
}

func NewDeleteRequest(conf config.Config, graphName, transaction string) DeleteRequest {
	return DeleteRequest{request: request{
		config:      conf,
		graphName:   graphName,
		transaction: transaction,
		operation:   "D",
	},
	}
}

func (r DeleteRequest) Send(filename string, modified time.Time) error {
	url := fmt.Sprintf("%s/%s/delete/%s?ta_id=%s&modified_date=%d", r.config.Server.Host, r.graphName, filename, r.transaction, modified.UnixMilli())
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

func (u UploadRequest) Send(filename string, modified time.Time, body []byte) error {
	return u.upload(filename, modified, body)
}

func (r request) upload(filename string, modified time.Time, body []byte) error {
	url := fmt.Sprintf("%s/%s/upload", r.config.Server.Host, r.graphName)
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	fileWriter, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = fileWriter.Write(body)
	if err != nil {
		return err
	}

	err = addFormField(mw, "ta-id", r.transaction)
	if err != nil {
		return err
	}

	err = addFormField(mw, "operation", r.operation)
	if err != nil {
		return err
	}

	err = addFormField(mw, "modified-date", modified.Format(time.RFC3339))
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
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Message: %s\n", respBody)
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
