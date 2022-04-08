package repository

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type localFileSystemRepository struct {
	rootPath string
}

func NewLocalFileSystemRepository() *localFileSystemRepository {
	return &localFileSystemRepository{rootPath: os.Getenv("FILE_ROOT_PATH")}
}

func (r *localFileSystemRepository) StoreFile(fileName string, byteData []byte) error {
	f, err := os.OpenFile(r.rootPath+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Error during opening file: %v", err)
	}

	_, err = io.Copy(f, bytes.NewReader(byteData))
	if err != nil {
		return fmt.Errorf("Error during writing in file: %v", err)
	}

	return f.Close()
}

func (r *localFileSystemRepository) FindFileByPattern(pattern string) (*os.File, error) {
	matches, err := filepath.Glob(r.rootPath + pattern)
	if err != nil {
		return nil, err
	}

	if len(matches) != 1 {
		return nil, fmt.Errorf("Unable unique identify file with name: %s", pattern)
	}

	return os.Open("./" + matches[0])
}
