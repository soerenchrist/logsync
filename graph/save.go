package graph

import (
	"encoding/json"
	"io"
)

func SaveGraph(graph Graph, outputFile io.Writer) error {
	jsonString, err := json.Marshal(graph)
	if err != nil {
		return err
	}

	_, err = outputFile.Write(jsonString)
	if err != nil {
		return err
	}

	return nil
}
