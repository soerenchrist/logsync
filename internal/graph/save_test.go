package graph

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestLoadGraph(t *testing.T) {
	graph := getTestGraph()
	jsonData, _ := json.Marshal(graph)

	buffer := bytes.NewBuffer(jsonData)

	g, err := LoadGraph(buffer)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}

	if g.Name != "test" {
		t.Fatalf("expected name \"test\", got %s", g.Name)
	}

	if len(g.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(g.Files))
	}
}

func TestSaveGraph(t *testing.T) {
	graph := getTestGraph()

	writer := bytes.Buffer{}

	err := SaveGraph(graph, &writer)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}

	result := writer.String()
	if result != "{\"name\":\"test\",\"files\":[{\"id\":\"Id1\",\"path\":\"test/id1\",\"lastChange\":\"2024-03-18T19:24:59.418+01:00\"},{\"id\":\"Id2\",\"path\":\"test/id2\",\"lastChange\":\"2024-03-18T19:24:59.418+01:00\"}]}" {
		t.Fatalf("Got wrong json: %s", result)
	}
}

func getTestGraph() Graph {
	return Graph{
		Name: "test",
		Files: []File{
			{
				Id:         "Id1",
				Path:       "test/id1",
				LastChange: time.UnixMilli(1710786299418),
			},
			{
				Id:         "Id2",
				Path:       "test/id2",
				LastChange: time.UnixMilli(1710786299418),
			},
		},
	}
}
