package sync

import (
	"github.com/google/uuid"
	"github.com/soerenchrist/logsync/client/internal/compare"
	"github.com/soerenchrist/logsync/client/internal/config"
	"github.com/soerenchrist/logsync/client/internal/crypt"
	"github.com/soerenchrist/logsync/client/internal/graph"
	"github.com/soerenchrist/logsync/client/internal/log"
	"github.com/soerenchrist/logsync/client/internal/remote"
	"io"
	"os"
	"slices"
	"time"
)

type graphSyncer struct {
	config      config.Config
	remote      *remote.Remote
	savedGraph  *graph.Graph
	basePath    string
	transaction string
	name        string
}

func newSyncer(graphPath string, conf config.Config) (graphSyncer, error) {
	r := remote.New(conf)
	name, err := graph.GetNameByPath(graphPath)
	if err != nil {
		return graphSyncer{}, err
	}
	transaction, _ := uuid.NewUUID()
	log.Info("Graph name: %s", name)

	loadFilePath, err := getLoadFilePath(name)
	if err != nil {
		return graphSyncer{}, err
	}
	savedGraph, err := graph.LoadGraphFromFile(loadFilePath)
	if err != nil {
		return graphSyncer{}, err
	}
	return graphSyncer{
		config:      conf,
		transaction: transaction.String(),
		basePath:    graphPath,
		savedGraph:  &savedGraph,
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
	for _, graphPath := range conf.Sync.Graphs {
		syncer, err := newSyncer(graphPath, conf)
		if err != nil {
			log.Error("Could not create syncer", err)
			continue
		}
		err = syncer.syncGraph()
		if err != nil {
			log.Error("Failed to sync", err)
		}
	}
}

func (s graphSyncer) syncGraph() error {
	log.Info("Last sync was %v", s.savedGraph.LastSync)

	remoteChanges, err := s.remote.GetChanges(s.name, s.savedGraph.LastSync)
	if err != nil {
		return err
	}
	log.Info("Found %d remote changes", len(remoteChanges))

	readGraph, err := graph.ReadGraph(s.basePath)
	if err != nil {
		return err
	}

	localChanges, err := s.getLocalChanges(readGraph)
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

	s.savedGraph.LastSync = time.Now()
	err = graph.SaveGraphToFile(*s.savedGraph, savePath)
	if err != nil {
		return err
	}

	return nil
}

func (s graphSyncer) deleteFile(file graph.File) error {
	var fileId string
	var err error
	if s.config.Encryption.Enabled {
		log.Info("Encrypting filename")
		fileId, err = crypt.EncryptString(file.Id, s.config.Encryption.Key)
		if err != nil {
			return err
		}
	} else {
		fileId = file.Id
	}

	request := remote.NewRequest(s.config, s.name, s.transaction, "D")
	return request.SendDelete(fileId, file.LastChange)
}

func (s graphSyncer) uploadFile(file graph.File, operation string) error {
	request := remote.NewRequest(s.config, s.name, s.transaction, operation)

	f, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var body []byte
	var fileId string
	if s.config.Encryption.Enabled {
		log.Info("Encrypting content")
		body, err = crypt.Encrypt(contents, s.config.Encryption.Key)
		if err != nil {
			return err
		}
		fileId, err = crypt.EncryptString(file.Id, s.config.Encryption.Key)
		if err != nil {
			return err
		}
	} else {
		body = contents
		fileId = file.Id
	}

	return request.SendUpload(fileId, file.LastChange, body)
}

func (s graphSyncer) uploadChanges(changes compare.Result, conflicts []string) error {
	log.Info("Uploading changes to server")
	for _, created := range changes.Created {
		if slices.Contains(conflicts, created.Id) {
			log.Info("Skipping upload for conflict file %s", created.Id)
			continue
		}
		log.Info("Uploading created file: %s", created.Id)
		err := s.uploadFile(created, "C")
		if err != nil {
			log.Error("Failed to upload", err)
			continue
		}
		s.savedGraph.AddOrUpdateFile(created)
	}

	for _, changed := range changes.Changed {
		if slices.Contains(conflicts, changed.Id) {
			log.Info("Skipping upload for conflict file %s", changed.Id)
			continue
		}
		log.Info("Uploading changed file: %s", changed.Id)
		err := s.uploadFile(changed, "M")
		if err != nil {
			log.Error("Failed to upload change", err)
		}
		s.savedGraph.AddOrUpdateFile(changed)
	}

	for _, deleted := range changes.Deleted {
		if slices.Contains(conflicts, deleted.Id) {
			log.Info("Skipping deletion for conflict file %s", deleted.Id)
			continue
		}
		log.Info("Deleting file: %s", deleted.Id)
		err := s.deleteFile(deleted)
		if err != nil {
			log.Error("Failed to delete", err)
		}
		s.savedGraph.RemoveFile(deleted.Id)
	}

	return nil
}

func (s graphSyncer) downloadFile(fileId string) error {
	content, err := s.remote.GetContent(s.name, fileId)
	if err != nil {
		log.Error("Failed to download content", err)
		return err
	}
	if s.config.Encryption.Enabled {
		content, err = crypt.Decrypt(content, s.config.Encryption.Key)
		if err != nil {
			return err
		}

		fileId, err = crypt.DecryptString(fileId, s.config.Encryption.Key)
		if err != nil {
			return err
		}
	}
	path, err := graph.StoreFile(s.basePath, fileId, content)
	if err != nil {
		log.Error("Failed to store file in local graph", err)
		return err
	}
	s.savedGraph.AddOrUpdateFile(graph.File{
		Id:         fileId,
		Path:       path,
		LastChange: time.Now(),
	})

	return nil
}

func (s graphSyncer) removeFile(fileId string) error {
	var err error
	if s.config.Encryption.Enabled {
		fileId, err = crypt.DecryptString(fileId, s.config.Encryption.Key)
		if err != nil {
			return err
		}
	}
	err = graph.RemoveFile(s.basePath, fileId)
	if err != nil {
		return err
	}
	s.savedGraph.RemoveFile(fileId)
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
			err := s.downloadFile(change.FileId)
			if err != nil {
				log.Error("Failed to store file in local graph", err)
				continue
			}
		} else if change.Operation == "D" {
			err := s.removeFile(change.FileId)
			if err != nil {
				log.Error("Failed to remove file in local graph", err)
				continue
			}
		}
	}
	return nil
}

func (s graphSyncer) getLocalChanges(g graph.Graph) (compare.Result, error) {
	compResult := compare.Graphs(*s.savedGraph, g)
	return compResult, nil
}
