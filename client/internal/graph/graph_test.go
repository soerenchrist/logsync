package graph

import (
	"errors"
	"io/fs"
	"testing"
)

func TestReadGraph(t *testing.T) {
	t.Run("graph exists", func(t *testing.T) {
		graph, err := ReadGraph("testdata/graph")
		if err != nil {
			t.Fatalf("Should not fail with err: %v", err)
		}

		if graph.Name != "graph" {
			t.Fatalf("Expected Name %s, got %s", "graph", graph.Name)
		}

		if len(graph.Files) != 4 {
			t.Fatalf("Should have length 4, has %d", len(graph.Files))
		}

		tt := []struct {
			id   string
			path string
		}{
			{
				id:   "journals|2024_03_02.md",
				path: "testdata/graph/journals/2024_03_02.md",
			},
			{
				id:   "journals|2024_03_03.md",
				path: "testdata/graph/journals/2024_03_03.md",
			},
			{
				id:   "pages|Page1.md",
				path: "testdata/graph/pages/Page1.md",
			},
			{
				id:   "pages|Page2.md",
				path: "testdata/graph/pages/Page2.md",
			},
		}

		for index, found := range graph.Files {
			expected := tt[index]
			if expected.id != found.Id {
				t.Fatalf("Expected Id %s, got %s", expected.id, found.Id)
			}

			if expected.path != found.Path {
				t.Fatalf("Expected Path %s, got %s", expected.path, found.Path)
			}
		}
	})

	t.Run("dir does not exist", func(t *testing.T) {
		_, err := ReadGraph("testdata/doesNotExist")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("expected ErrNotExist, got %v", err)
		}
	})
}
