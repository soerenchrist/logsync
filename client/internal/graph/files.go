package graph

import (
	"errors"
	"os"
	"path"
	"strings"
)

func StoreFile(graphPath, fileId string, content []byte) error {
	p := getPathByFileId(fileId)
	p = path.Join(graphPath, p)

	err := ensureDirExists(p)
	if err != nil {
		return err
	}

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

func ensureDirExists(p string) error {
	dir := path.Dir(p)
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}

// RemoveFile TODO: maybe introduce some kind of trash bin
func RemoveFile(graphPath, fileId string) error {
	p := getPathByFileId(fileId)
	p = path.Join(graphPath, p)

	return os.Remove(p)
}

func getPathByFileId(fileId string) string {
	parts := strings.Split(fileId, Separator)
	return path.Join(parts...)
}
