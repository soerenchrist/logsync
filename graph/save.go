package graph

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveGraph(graph Graph, outputFile string) error {
	jsonString, err := json.Marshal(graph)
	if err != nil {
		return err
	}

	fmt.Printf(string(jsonString))

	file, err := os.Create(outputFile)
	defer file.Close()

	if err != nil {
		return err
	}
	_, err = file.Write(jsonString)
	if err != nil {
		return err
	}

	return nil
}
