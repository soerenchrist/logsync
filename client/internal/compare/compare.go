package compare

import (
	"github.com/soerenchrist/logsync/client/internal/graph"
	"slices"
)

type CompResult struct {
	changed []graph.File
	created []graph.File
	deleted []graph.File
}

func (c CompResult) NoChanges() bool {
	return len(c.changed) == 0 && len(c.deleted) == 0 && len(c.created) == 0
}

func Graphs(old graph.Graph, new graph.Graph) CompResult {
	created := make([]graph.File, 0)
	changed := make([]graph.File, 0)
	deleted := make([]graph.File, 0)
	for _, newFile := range new.Files {
		foundOld := find(old.Files, newFile.Id)
		if foundOld == nil {
			created = append(created, newFile)
		} else {
			if newFile.LastChange.After(foundOld.LastChange) {
				changed = append(changed, newFile)
			}
		}
	}

	for _, oldFile := range old.Files {
		foundNew := find(new.Files, oldFile.Id)
		if foundNew == nil {
			deleted = append(deleted, oldFile)
		}
	}

	return CompResult{
		created: created,
		deleted: deleted,
		changed: changed,
	}
}

func find(files []graph.File, id string) *graph.File {
	index := slices.IndexFunc(files, func(file graph.File) bool {
		return file.Id == id
	})

	if index < 0 {
		return nil
	}

	return &files[index]
}
