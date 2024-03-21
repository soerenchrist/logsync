package graph

import (
	"os"
	"path"
	"strings"
)

func StoreFile(fileId string, content []byte) error {
	p := getPathByFileId(fileId)

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

// RemoveFile TODO: maybe introduce some kind of trash bin
func RemoveFile(fileId string) error {
	p := getPathByFileId(fileId)

	return os.Remove(p)
}

func getPathByFileId(fileId string) string {
	parts := strings.Split(fileId, Separator)
	return path.Join(parts...)
}
