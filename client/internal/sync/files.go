package sync

import (
	"github.com/soerenchrist/logsync/client/internal/graph"
	"os"
	"path"
	"strings"
)

func storeFile(fileId string, content []byte) error {
	p := getPathByFileId(fileId)

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

func getPathByFileId(fileId string) string {
	parts := strings.Split(fileId, graph.Separator)
	return path.Join(parts...)
}
