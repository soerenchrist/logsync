package sync

import (
	"errors"
	"os"
	"path"
	"time"
)

func saveLastSyncTime(timestamp time.Time) error {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := path.Join(dirName, ".config", "logsync", ".lastsync")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}

func getLastSyncTime() (time.Time, error) {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return time.Time{}, err
	}

	filePath := path.Join(dirName, ".config", "logsync", ".lastsync")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, string(data))
}
