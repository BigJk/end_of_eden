//go:build js
// +build js

package fs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall/js"
)

type noOpWriteCloser struct{}

func (noOpWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (noOpWriteCloser) Close() error {
	return nil
}

var fileIndex = make(map[string]FileInfo)

func init() {
	data, err := ReadFile("/assets/file_index.json")
	if err != nil {
		panic(err)
	}

	var fis []FileInfo
	if err := json.Unmarshal(data, &fis); err != nil {
		panic(err)
	}
	for _, fi := range fis {
		fi.Path = filepath.Clean(fi.Path)
		fileIndex[fi.Path] = fi
	}

	js.Global().Set("fsDump", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		for _, fi := range fis {
			fmt.Println(fi.Path)
		}
		return nil
	}))
}

func ReadDir(path string) ([]FileInfo, error) {
	cleanPath := filepath.Clean(path)
	var fis []FileInfo
	for indexPath := range fileIndex {
		if strings.HasPrefix(indexPath, cleanPath) {
			fis = append(fis, fileIndex[path])
		}
	}
	return fis, nil
}

func ReadFile(path string) ([]byte, error) {
	// Check for temp file
	jsRes := js.Global().Call("fsRead", path)
	if !jsRes.IsNull() && !jsRes.IsUndefined() {
		return base64.StdEncoding.DecodeString(jsRes.String())
	}

	// Check for asset
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not load file %s: %s", path, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func OpenFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	// TODO: Implement
	return noOpWriteCloser{}, nil
}

func WriteFile(path string, data []byte) error {
	// TODO: error handling
	_ = js.Global().Call("fsWrite", path, base64.StdEncoding.EncodeToString(data))
	return nil
}

func Walk(root string, walkFn func(path string, isDir bool) error) error {
	keys := lo.Keys(fileIndex)
	sort.Strings(keys)

	cleanPath := filepath.Clean(root)
	for _, path := range keys {
		if !strings.HasPrefix(path, cleanPath) {
			continue
		}

		if err := walkFn(path, fileIndex[path].IsDir()); err != nil {
			return err
		}
	}
	return nil
}
