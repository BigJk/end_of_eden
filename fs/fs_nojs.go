//go:build !js
// +build !js

package fs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func OpenFile(path string, flag int, perm os.FileMode) (io.ReadSeekCloser, error) {
	return os.OpenFile(path, flag, perm)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ReadDir(path string) ([]FileInfo, error) {
	res, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var files []FileInfo
	for _, f := range res {
		files = append(files, FileInfo{Name: f.Name(), IsDir: f.IsDir()})
	}
	return files, nil
}

func Walk(path string, walkFn func(path string, info FileInfo, err error) error) error {
	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		return walkFn(path, FileInfo{Name: info.Name(), IsDir: info.IsDir()}, err)
	})
}

func WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}
