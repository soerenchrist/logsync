package graph

import (
	"errors"
	"github.com/soerenchrist/logsync/client/internal/sanitize"
	"os"
	"path"
	"slices"
	"time"
)

var Separator = "___"

var skipFolders = []string{"logseq"}

type Graph struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
}

type File struct {
	Id         string    `json:"id"`
	Path       string    `json:"path"`
	LastChange time.Time `json:"lastChange"`
}

func ReadGraph(baseDir string) (Graph, error) {
	files := make([]File, 0)
	errs := make([]error, 0)
	traverseGraph(baseDir, "", &files, &errs)

	graphName, err := getGraphName(baseDir)
	if err != nil {
		return Graph{}, err
	}

	if len(errs) > 0 {
		return Graph{}, errors.Join(errs...)
	}

	return Graph{
		Files: files,
		Name:  graphName,
	}, nil
}

func getGraphName(baseDir string) (string, error) {
	stat, err := os.Stat(baseDir)
	if err != nil {
		return "", err
	}
	return stat.Name(), nil
}

func traverseGraph(baseDir string, name string, files *[]File, errors *[]error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	for _, entry := range entries {
		fileId := buildFileId(name, entry.Name())
		filePath := path.Join(baseDir, entry.Name())
		if entry.IsDir() {
			if slices.Contains(skipFolders, entry.Name()) {
				continue
			}
			traverseGraph(filePath, fileId, files, errors)
		} else {
			info, err := entry.Info()
			if err != nil {
				*errors = append(*errors, err)
				return
			}
			file := File{
				Id:         fileId,
				Path:       sanitize.Path(filePath),
				LastChange: info.ModTime(),
			}
			*files = append(*files, file)
		}
	}
}

func buildFileId(baseName, name string) string {
	if len(baseName) == 0 {
		return name
	}

	return baseName + Separator + name
}

func GetNameByPath(graphPath string) (string, error) {
	stat, err := os.Stat(graphPath)
	if err != nil {
		return "", err
	}

	return stat.Name(), nil
}
