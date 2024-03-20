package graph

import (
	"encoding/json"
	"errors"
	"io"
	"os"
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

func SaveGraphToFile(graph Graph, outputFile string) error {
	file, err := os.Create(outputFile)
	defer file.Close()
	if err != nil {
		return err
	}

	return SaveGraph(graph, file)
}

func LoadGraph(inputFile io.Reader) (Graph, error) {
	var g Graph

	decoder := json.NewDecoder(inputFile)
	err := decoder.Decode(&g)
	if err != nil {
		return Graph{}, err
	}

	return g, nil
}

func LoadGraphFromFile(filePath string) (Graph, error) {
	file, err := os.Open(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return Graph{}, nil
	}

	return LoadGraph(file)
}
