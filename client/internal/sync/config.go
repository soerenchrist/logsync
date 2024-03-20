package sync

import (
	"errors"
	"os"
	"path"
)

func getLoadFilePath(graphName string) (string, error) {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := path.Join(dirName, ".config")
	err = ensureCreated(configDir)
	if err != nil {
		return "", err
	}

	logsyncDir := path.Join(configDir, "logsync")
	err = ensureCreated(logsyncDir)
	if err != nil {
		return "", err
	}

	return path.Join(dirName, ".config", "logsync", graphName+".json"), nil
}

func ensureCreated(dir string) error {
	_, err := os.Stat(dir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
