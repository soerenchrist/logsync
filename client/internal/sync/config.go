package sync

import (
	"os"
	"path"
)

func getLoadFilePath(graphName string) (string, error) {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(dirName, ".config", "logsync", graphName+".json"), nil
}
