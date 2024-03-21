package graph

import (
	"os"
	"testing"
)

func TestStoreFile(t *testing.T) {
	t.Run("dir already exists", func(t *testing.T) {
		err := StoreFile("testdata/graph", "journals___stored.md", []byte{1, 2, 3})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		_, err = os.Stat("testdata/graph/journals/stored.md")
		if err != nil {
			t.Fatalf("File should exist, expected no error, got %v", err)
		}

		_ = os.Remove("testdata/graph/journals/stored.md")
	})

	t.Run("missing dir", func(t *testing.T) {
		err := StoreFile("testdata/graph", "something___stored.md", []byte{1, 2, 3})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		_, err = os.Stat("testdata/graph/something/stored.md")
		if err != nil {
			t.Fatalf("File should exist, expected no error, got %v", err)
		}

		_ = os.Remove("testdata/graph/something/stored.md")
		_ = os.Remove("testdata/graph/something")
	})
}
