package sync

import (
	"fmt"
	"github.com/soerenchrist/logsync/client/internal/compare"
	"github.com/soerenchrist/logsync/client/internal/config"
	"github.com/soerenchrist/logsync/client/internal/graph"
	"github.com/soerenchrist/logsync/client/internal/remote"
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
		fmt.Printf("Remote change %d: %v\n", i, change)
	}

	localChanges, err := getLocalChanges(graphPath, name)

	conflicts := checkForConflicts(remoteChanges, localChanges)
	for _, conflict := range conflicts {
		fmt.Printf("Conflict found for file: %s\n", conflict)
	}

	err = downloadChanges(remoteChanges, conflicts)
	if err != nil {
		return err
	}
	err = uploadChanges(name, localChanges, conflicts)
	if err != nil {
		return err
	}

	return nil
}

func uploadChanges(graphName string, changes compare.Result, conflicts []string) error {
	fmt.Printf("Uploading created files")
	for _, created := range changes.Created {
		err := r.UploadFile(graphName, created, "C")
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}

	return nil
}

func downloadChanges(changes []remote.ChangeLogEntry, conflicts []string) error {

	return nil
}

func getLocalChanges(graphPath string, graphName string) (compare.Result, error) {
	readGraph, err := graph.ReadGraph(graphPath)
	if err != nil {
		return compare.Result{}, err
	}

	loadFilePath, err := getLoadFilePath(graphName)
	if err != nil {
		return compare.Result{}, err
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
			fmt.Printf("Failed to sync: %v\n", err)
		}
	}
	err := saveLastSyncTime(time.Now())
	if err != nil {
		fmt.Printf("Failed to save last sync time: %v\n", err)
	}
}
