package sync

import (
	"github.com/google/uuid"
	"github.com/soerenchrist/logsync/client/internal/compare"
	"github.com/soerenchrist/logsync/client/internal/config"
	"github.com/soerenchrist/logsync/client/internal/graph"
	"github.com/soerenchrist/logsync/client/internal/log"
	"github.com/soerenchrist/logsync/client/internal/remote"
	"slices"
	"time"
)

var r *remote.Remote

func Start(conf config.Config) {
	r = remote.New(conf)
	ticker := time.Tick(time.Duration(conf.Sync.Interval) * time.Second)
	for range ticker {
		log.Info("Starting sync of graphs")
		syncGraphs(conf.Sync.Graphs)
	}
}

func syncGraph(graphPath string) error {
	name, err := graph.GetNameByPath(graphPath)
	if err != nil {
		return err
	}
	transaction, _ := uuid.NewUUID()
	log.Info("Graph name: %s", name)

	lastSync, err := getLastSyncTime()
	if err != nil {
		return err
	}
	log.Info("Last sync was %v", lastSync)

	remoteChanges, err := r.GetChanges(name, lastSync)
	if err != nil {
		return err
	}
	log.Info("Found %d remote changes", len(remoteChanges))

	readGraph, err := graph.ReadGraph(graphPath)
	if err != nil {
		return err
	}

	localChanges, err := getLocalChanges(readGraph)
	if err != nil {
		return err
	}

	conflicts := checkForConflicts(remoteChanges, localChanges)
	log.Info("Found %d conflicts", len(conflicts))

	err = downloadChanges(remoteChanges, conflicts)
	if err != nil {
		return err
	}
	err = uploadChanges(name, localChanges, transaction, conflicts)
	if err != nil {
		return err
	}

	savePath, err := getLoadFilePath(readGraph.Name)
	if err != nil {
		return err
	}
	err = graph.SaveGraphToFile(readGraph, savePath)
	if err != nil {
		return err
	}

	return nil
}

func uploadChanges(graphName string, changes compare.Result, transaction uuid.UUID, conflicts []string) error {
	log.Info("Uploading changes to server")
	for _, created := range changes.Created {
		if slices.Contains(conflicts, created.Id) {
			log.Info("Skipping upload for conflict file %s", created.Id)
			continue
		}
		log.Info("Uploading created file: %s", created.Id)
		err := r.UploadFile(graphName, created, transaction.String(), "C")
		if err != nil {
			log.Error("Failed to upload", err)
		}
	}

	for _, changed := range changes.Changed {
		if slices.Contains(conflicts, changed.Id) {
			log.Info("Skipping upload for conflict file %s", changed.Id)
			continue
		}
		log.Info("Uploading changed file: %s", changed.Id)
		err := r.UploadFile(graphName, changed, transaction.String(), "M")
		if err != nil {
			log.Error("Failed to upload change", err)
		}
	}

	for _, deleted := range changes.Deleted {
		if slices.Contains(conflicts, deleted.Id) {
			log.Info("Skipping deletion for conflict file %s", deleted.Id)
			continue
		}
		log.Info("Deleting file: %s", deleted.Id)
		err := r.DeleteFile(graphName, deleted, transaction.String())
		if err != nil {
			log.Error("Failed to delete", err)
		}
	}

	return nil
}

func downloadChanges(changes []remote.ChangeLogEntry, conflicts []string) error {
	log.Info("Downloading changes from server")
	for _, change := range changes {
		if slices.Contains(conflicts, change.FileId) {
			log.Info("Skipping download of file %s", change.FileId)
			continue
		}
		log.Info("Found change with transaction %s for file %s", change.Transaction, change.FileId)
		if change.Operation == "C" || change.Operation == "M" {
			content, err := r.GetContent(change.GraphName, change.FileId)
			if err != nil {
				log.Error("Failed to download content", err)
				continue
			}
			err = storeFile(change.FileId, content)
			if err != nil {
				log.Error("Failed to store file in local graph", err)
				continue
			}
		}
	}
	return nil
}

func getLocalChanges(g graph.Graph) (compare.Result, error) {
	loadFilePath, err := getLoadFilePath(g.Name)
	if err != nil {
		return compare.Result{}, err
	}
	log.Info("Save file path: %s", loadFilePath)
	savedGraph, err := graph.LoadGraphFromFile(loadFilePath)
	if err != nil {
		return compare.Result{}, err
	}

	compResult := compare.Graphs(savedGraph, g)
	return compResult, nil
}

func syncGraphs(graphs []string) {
	for _, graphPath := range graphs {
		err := syncGraph(graphPath)
		if err != nil {
			log.Error("Failed to sync", err)
		}
	}
	err := saveLastSyncTime(time.Now())
	if err != nil {
		log.Error("Failed to save last sync time", err)
	}
}
