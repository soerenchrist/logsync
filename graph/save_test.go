package graph

import (
	"bytes"
	"testing"
	"time"
)

func TestSaveGraph(t *testing.T) {
	graph := Graph{
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
