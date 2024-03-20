package remote

import (
	"encoding/json"
	"fmt"
	"github.com/soerenchrist/logsync/client/internal/config"
	"io"
	"net/http"
	"time"
)

type Remote struct {
	conf config.Config
}

type ChangeLogEntry struct {
	GraphName string    `json:"graph_name"`
	FileId    string    `json:"file_id"`
	Timestamp time.Time `json:"timestamp"`
	FileName  string    `json:"file_name"`
	Operation string    `json:"operation"`
}

func New(conf config.Config) *Remote {
	return &Remote{conf: conf}
}

func (r *Remote) GetChanges(graphName string, since time.Time) ([]ChangeLogEntry, error) {
	url := fmt.Sprintf("%s/%s/changes?since=%s", r.conf.Server.Host, graphName, since.Format(time.RFC3339))
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var entries []ChangeLogEntry
	err = json.Unmarshal(body, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
