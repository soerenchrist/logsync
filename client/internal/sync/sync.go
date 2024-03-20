package sync

import (
	"fmt"
	"github.com/soerenchrist/logsync/client/internal/compare"
	"github.com/soerenchrist/logsync/client/internal/config"
	"github.com/soerenchrist/logsync/client/internal/graph"
	"github.com/soerenchrist/logsync/client/internal/remote"
	"os"
	"path"
	"time"
)

var r *remote.Remote

func Start(conf config.Config) {
	r = remote.New(conf)
	ticker := time.NewTicker(time.Duration(conf.Sync.Interval) * time.Second)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		syncGraphs(conf.Sync.Graphs)
		fmt.Println("Tick")
	}
}

func syncGraph(graphPath string) error {
	name, err := graph.GetNameByPath(graphPath)
	if err != nil {
		return err
	}
	fmt.Printf("Graph name: %s\n", name)

	lastSync, err := getLastSyncTime()
	if err != nil {
		return err
	}

	remoteChanges, err := r.GetChanges(name, lastSync)
	if err != nil {
		return err
	}

	for i, change := range remoteChanges {
		fmt.Printf("%d: %v", i, change)
	}

	_, err = getLocalChanges(graphPath, name)

	return nil

	/*
		loadFilePath, err := getLoadFilePath(name)
		if err != nil {
			return err
		}
		fmt.Printf("Save file path: %s\n", loadFilePath)
		savedGraph, err := graph.LoadGraphFromFile(loadFilePath)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Local graph does not exist -> loading it for the first time\n")
			return firstLoad(graphPath, loadFilePath)
		} else if err != nil {
			return err
		}

		readGraph, err := graph.ReadGraph(graphPath)
		if err != nil {
			return err
		}

		compResult := compare.Graphs(savedGraph, readGraph)
		if compResult.NoChanges() {
			fmt.Printf("No remoteChanges\n")
		} else {
			fmt.Printf("Changes\n")
		}

		return nil
	*/
}

func getLocalChanges(graphPath string, graphName string) (compare.CompResult, error) {
	readGraph, err := graph.ReadGraph(graphPath)
	if err != nil {
		return compare.CompResult{}, err
	}

	loadFilePath, err := getLoadFilePath(graphName)
	if err != nil {
		return compare.CompResult{}, err
	}
	fmt.Printf("Save file path: %s\n", loadFilePath)
	savedGraph, err := graph.LoadGraphFromFile(loadFilePath)

	compResult := compare.Graphs(savedGraph, readGraph)
	if compResult.NoChanges() {
		fmt.Printf("No changes\n")
	} else {
		fmt.Printf("Changes\n")
	}

	return compResult, nil
}

func firstLoad(graphPath string, filePath string) error {
	g, err := graph.ReadGraph(graphPath)
	if err != nil {
		return err
	}

	return graph.SaveGraphToFile(g, filePath)
}

func syncGraphs(graphs []string) {
	for _, graphPath := range graphs {
		err := syncGraph(graphPath)
		if err != nil {
			fmt.Printf("Failed to sync: %v", err)
		}
	}
}

func getLoadFilePath(graphName string) (string, error) {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(dirName, ".config", "logsync", graphName+".json"), nil
}
