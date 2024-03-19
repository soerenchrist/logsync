package files

import (
	"errors"
	"io"
	"os"
	"path"
)

type FileStore interface {
	Store(graphName string, fileName string, reader io.Reader) error
	Remove(graphName string, fileName string) error
}

type Files struct {
	basePath string
}

func New(basePath string) Files {
	return Files{
		basePath: basePath,
	}
}

func (f Files) Store(graphName string, fileName string, reader io.Reader) error {
	err := f.ensureGraphDirExists(graphName)
	if err != nil {
		return err
	}

	filePath := path.Join(f.basePath, graphName, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func (f Files) Remove(graphName string, fileName string) error {
	filePath := path.Join(f.basePath, graphName, fileName)

	return os.Remove(filePath)
}

func (f Files) ensureGraphDirExists(graphName string) error {
	err := f.ensureExists(f.basePath)
	if err != nil {
		return err
	}

	graphPath := path.Join(f.basePath, graphName)
	return f.ensureExists(graphPath)
}

func (f Files) ensureExists(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
