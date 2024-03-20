package compare

import (
	"github.com/soerenchrist/logsync/client/internal/graph"
	"slices"
)

type Result struct {
	Changed []graph.File
	Created []graph.File
	Deleted []graph.File
}

func (c Result) NoChanges() bool {
	return len(c.Changed) == 0 && len(c.Deleted) == 0 && len(c.Created) == 0
}

func Graphs(old graph.Graph, new graph.Graph) Result {
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

	return Result{
		Created: created,
		Deleted: deleted,
		Changed: changed,
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
