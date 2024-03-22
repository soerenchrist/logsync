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

type graphSyncer struct {
	remote      *remote.Remote
	basePath    string
	transaction string
	name        string
}

func newSyncer(graphPath string, r *remote.Remote) (graphSyncer, error) {
	name, err := graph.GetNameByPath(graphPath)
	if err != nil {
		return graphSyncer{}, err
	}
	transaction, _ := uuid.NewUUID()
	log.Info("Graph name: %s", name)
	return graphSyncer{
		transaction: transaction.String(),
		basePath:    graphPath,
		remote:      r,
		name:        name,
	}, nil
}

func Start(conf config.Config) {
	if conf.Sync.Once {
		log.Info("Syncing graphs once")
		syncGraphs(conf)
		return
	}

	ticker := time.Tick(time.Duration(conf.Sync.Interval) * time.Second)
	for range ticker {
		log.Info("Starting sync of graphs")
		syncGraphs(conf)
	}
}

func syncGraphs(conf config.Config) {
	r := remote.New(conf)
	for _, graphPath := range conf.Sync.Graphs {
		syncer, err := newSyncer(graphPath, r)
		if err != nil {
			log.Error("Could not create syncer", err)
			continue
		}
		err = syncer.syncGraph()
		if err != nil {
			log.Error("Failed to sync", err)
		}
	}
	err := saveLastSyncTime(time.Now())
	if err != nil {
		log.Error("Failed to save last sync time", err)
	}
}

func (s graphSyncer) syncGraph() error {
	lastSync, err := getLastSyncTime()
	if err != nil {
		return err
	}
	log.Info("Last sync was %v", lastSync)

	remoteChanges, err := s.remote.GetChanges(s.name, lastSync)
	if err != nil {
		return err
	}
	log.Info("Found %d remote changes", len(remoteChanges))

	readGraph, err := graph.ReadGraph(s.basePath)
	if err != nil {
		return err
	}

	localChanges, err := getLocalChanges(readGraph)
	if err != nil {
		return err
	}

	// TODO: handle conflicts
	conflicts := checkForConflicts(remoteChanges, localChanges)
	log.Info("Found %d conflicts", len(conflicts))

	err = s.downloadChanges(remoteChanges, conflicts)
	if err != nil {
		return err
	}
	err = s.uploadChanges(localChanges, conflicts)
	if err != nil {
		return err
	}

	savePath, err := getLoadFilePath(readGraph.Name)
	if err != nil {
		return err
	}

	// need to reread the graph after updating
	// TODO: could do the update in memory to save the second read
	readGraph, _ = graph.ReadGraph(s.name)
	err = graph.SaveGraphToFile(readGraph, savePath)
	if err != nil {
		return err
	}

	return nil
}

func (s graphSyncer) uploadChanges(changes compare.Result, conflicts []string) error {
	log.Info("Uploading changes to server")
	for _, created := range changes.Created {
		if slices.Contains(conflicts, created.Id) {
			log.Info("Skipping upload for conflict file %s", created.Id)
			continue
		}
		log.Info("Uploading created file: %s", created.Id)
		err := s.remote.UploadFile(s.name, created, s.transaction, "C")
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
		err := s.remote.UploadFile(s.name, changed, s.transaction, "M")
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
		err := s.remote.DeleteFile(s.name, deleted, s.transaction)
		if err != nil {
			log.Error("Failed to delete", err)
		}
	}

	return nil
}

func (s graphSyncer) downloadChanges(changes []remote.ChangeLogEntry, conflicts []string) error {
	log.Info("Downloading changes from server")
	for _, change := range changes {
		if slices.Contains(conflicts, change.FileId) {
			log.Info("Skipping download of file %s", change.FileId)
			continue
		}
		log.Info("Found change with transaction %s for file %s", change.Transaction, change.FileId)
		if change.Operation == "C" || change.Operation == "M" {
			content, err := s.remote.GetContent(change.GraphName, change.FileId)
			if err != nil {
				log.Error("Failed to download content", err)
				continue
			}
			err = graph.StoreFile(s.basePath, change.FileId, content)
			if err != nil {
				log.Error("Failed to store file in local graph", err)
				continue
			}
		} else if change.Operation == "D" {
			err := graph.RemoveFile(s.basePath, change.FileId)
			if err != nil {
				log.Error("Failed to remove file in local graph", err)
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
