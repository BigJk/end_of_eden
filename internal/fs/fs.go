//go:build !js
// +build !js

package fs

import (
	"github.com/samber/lo"
	"io"
	"os"
	"path/filepath"
)

func ReadDir(path string) ([]FileInfo, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return lo.Map(dir, func(f os.DirEntry, i int) FileInfo {
		return FileInfo{
			Path:   filepath.Join(path, f.Name()),
			IsFile: !f.IsDir(),
		}
	}), nil
}

func OpenFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return os.OpenFile(name, flag, perm)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func Walk(root string, walkFn func(path string, isDir bool) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return walkFn(path, info.IsDir())
	})
}
