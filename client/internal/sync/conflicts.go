package sync

import (
	"github.com/soerenchrist/logsync/client/internal/compare"
	"github.com/soerenchrist/logsync/client/internal/graph"
	"github.com/soerenchrist/logsync/client/internal/remote"
	"slices"
)

func checkForConflicts(remoteChanges []remote.ChangeLogEntry, localChanges compare.Result) []string {
	if len(remoteChanges) == 0 {
		return []string{}
	}

	if localChanges.NoChanges() {
		return []string{}
	}

	conflicts := make([]string, 0)

	for _, remoteChange := range remoteChanges {
		if slices.ContainsFunc(localChanges.Changed, func(file graph.File) bool {
			return file.Id == remoteChange.FileId
		}) {
			conflicts = append(conflicts, remoteChange.FileId)
		}

		if slices.ContainsFunc(localChanges.Created, func(file graph.File) bool {
			return file.Id == remoteChange.FileId
		}) {
			conflicts = append(conflicts, remoteChange.FileId)
		}

		if slices.ContainsFunc(localChanges.Deleted, func(file graph.File) bool {
			return file.Id == remoteChange.FileId
		}) {
			conflicts = append(conflicts, remoteChange.FileId)
		}
	}

	return conflicts
}
