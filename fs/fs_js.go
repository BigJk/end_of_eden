//go:build js
// +build js

package fs

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

type file struct {
	*bytes.Reader
}

func (f *file) Close() error {
	return nil
}

type JsFS struct{}

func OpenFile(path string, flag int, perm os.FileMode) (io.ReadSeekCloser, error) {
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return &file{Reader: bytes.NewReader(body)}, nil
}

func ReadFile(path string) ([]byte, error) {
	f, err := OpenFile(path, 0, 0)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}

func ReadDir(path string) ([]FileInfo, error) {
	panic("implement me")
}

func Walk(path string, walkFn func(path string, info FileInfo, err error) error) error {
	panic("implement me")
}

func WriteFile(path string, data []byte, perm os.FileMode) error {
	panic("implement me")
}
