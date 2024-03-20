package routes

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/soerenchrist/logsync/server/internal/model"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpload(t *testing.T) {
	t.Run("No body", func(t *testing.T) {
		ts, _ := createTestServer()
		defer ts.Close()

		buf := bytes.Buffer{}

		req, _ := http.NewRequest("POST", ts.URL+"/test/upload", &buf)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected code 400, got %d", resp.StatusCode)
		}
	})

	t.Run("valid body", func(t *testing.T) {
		ts, db := createTestServer()
		defer ts.Close()

		req := createRequest("POST", ts.URL+"/test/upload", []byte{1, 2, 3})
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected code 201, got %d", resp.StatusCode)
		}

		var entries []model.ChangeLogEntry
		db.Find(&entries)

		if len(entries) != 1 {
			t.Fatalf("Expected 1 entry, got %d", len(entries))
		}
	})

}

func createTestServer() (*httptest.Server, *gorm.DB) {
	r := chi.NewRouter()
	db := createTestDb()
	c := NewController(db, r, TestStore{})
	c.MapEndpoints()

	return httptest.NewServer(r), db
}

func createTestDb() *gorm.DB {
	db, _ := model.CreateDb("file::memory:?cache=shared")
	return db
}

func createRequest(method, url string, file []byte) *http.Request {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	f, _ := mw.CreateFormFile("file", "file.txt")
	_, _ = f.Write(file)

	fileIdWriter, _ := mw.CreateFormField("file-id")
	_, _ = fileIdWriter.Write([]byte("testId"))

	opWriter, _ := mw.CreateFormField("operation")
	_, _ = opWriter.Write([]byte("C"))

	_ = mw.Close()

	req, _ := http.NewRequest(method, url, buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

type TestStore struct {
}

func (TestStore) Store(graphName string, fileName string, reader io.Reader) error {
	return nil
}

func (TestStore) Remove(graphName string, fileName string) error {
	return nil
}

func (TestStore) Content(graphName string, fileName string) ([]byte, error) {
	return nil, nil
}
