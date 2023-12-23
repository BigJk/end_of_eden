package fs

import "path/filepath"

type FileInfo struct {
	Path   string `json:"path"`
	IsFile bool   `json:"IsFile"`
}

func (fi FileInfo) Name() string {
	return filepath.Base(fi.Path)
}

func (fi FileInfo) IsDir() bool {
	return !fi.IsFile
}
