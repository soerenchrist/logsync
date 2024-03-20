package compare

import (
	"github.com/soerenchrist/logsync/client/internal/graph"
	"testing"
	"time"
)

func TestGraphs(t *testing.T) {
	t.Run("identical graphs", func(t *testing.T) {
		files := []graph.File{
			{
				Id:         "test1",
				Path:       "",
				LastChange: time.UnixMilli(100000000),
			},
			{
				Id:         "test2",
				Path:       "",
				LastChange: time.UnixMilli(100000000),
			},
		}
		graph1 := graph.Graph{
			Name:  "Graph1",
			Files: files,
		}

		graph2 := graph.Graph{
			Name:  "Graph1",
			Files: files,
		}

		res := Graphs(graph1, graph2)
		if len(res.Deleted) != 0 {
			t.Fatalf("Expected 0 deleted, got %d", len(res.Deleted))
		}
		if len(res.Changed) != 0 {
			t.Fatalf("Expected 0 changed, got %d", len(res.Changed))
		}
		if len(res.Created) != 0 {
			t.Fatalf("Expected 0 created, got %d", len(res.Created))
		}
	})

	t.Run("changed graphs", func(t *testing.T) {
		oldFiles := []graph.File{
			{
				Id:         "test1",
				Path:       "",
				LastChange: time.UnixMilli(100000000),
			},
			{
				Id:         "test2",
				Path:       "",
				LastChange: time.UnixMilli(100000000),
			},
			{
				Id:         "test3",
				Path:       "",
				LastChange: time.UnixMilli(100000000),
			},
		}
		newFiles := []graph.File{
			{
				Id:         "test1",
				Path:       "",
				LastChange: time.UnixMilli(100000000),
			},
			{
				Id:         "test2",
				Path:       "",
				LastChange: time.UnixMilli(200000000),
			},
			{
				Id:         "test4",
				Path:       "",
				LastChange: time.UnixMilli(100000000),
			},
		}
		oldGraph := graph.Graph{
			Name:  "Graph1",
			Files: oldFiles,
		}

		newGraph := graph.Graph{
			Name:  "Graph1",
			Files: newFiles,
		}

		res := Graphs(oldGraph, newGraph)
		if len(res.Deleted) != 1 {
			t.Fatalf("Expected 1 deleted, got %d", len(res.Deleted))
		}
		if len(res.Changed) != 1 {
			t.Fatalf("Expected 1 changed, got %d", len(res.Changed))
		}
		if len(res.Created) != 1 {
			t.Fatalf("Expected 1 created, got %d", len(res.Created))
		}
	})
}
