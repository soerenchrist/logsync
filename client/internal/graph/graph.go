package graph

import (
	"errors"
	"github.com/soerenchrist/logsync/client/internal/sanitize"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

var Separator = "___"

var skipFolders = []string{"bak", ".recycle"}

type Graph struct {
	Name     string    `json:"name"`
	LastSync time.Time `json:"lastSync"`
	Files    []File    `json:"files"`
}

func New(name string) Graph {
	return Graph{
		Name:     name,
		LastSync: time.Time{},
		Files:    make([]File, 0),
	}
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

func (g *Graph) AddOrUpdateFile(file File) {
	index := slices.IndexFunc(g.Files, func(file File) bool {
		return file.Id == file.Id
	})
	if index < 0 {
		g.Files = append(g.Files, file)
	} else {
		g.Files[index] = file
	}
}

func (g *Graph) RemoveFile(fileId string) {
	slices.DeleteFunc(g.Files, func(file File) bool {
		return file.Id == fileId
	})
}

func getGraphName(baseDir string) (string, error) {
	p := sanitize.Path(baseDir)
	parts := strings.Split(p, "/")
	if len(parts) < 1 {
		return "", errors.New("path is empty")
	}

	lastPart := parts[len(parts)-1]
	if strings.Contains(lastPart, ".") {
		parts := strings.Split(lastPart, ".")
		return parts[0], nil
	}
	return parts[len(parts)-1], nil
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
