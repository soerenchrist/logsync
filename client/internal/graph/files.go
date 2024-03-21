package graph

import (
	"os"
	"path"
	"strings"
)

func StoreFile(graphPath, fileId string, content []byte) error {
	p := getPathByFileId(fileId)
	p = path.Join(graphPath, p)

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
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
