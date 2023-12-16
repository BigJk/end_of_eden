package image

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	_ = os.Mkdir("./cache", 0755)
}

// hashFile returns a hash for the given file. If the file is not found in the
// search paths, an error is returned.
func hashFile(path string) (string, error) {
	for i := range searchPaths {
		data, err := ioutil.ReadFile(filepath.Join(searchPaths[i], path))
		if err != nil {
			continue
		}
		return fmt.Sprintf("%x", md5.Sum(data)), nil
	}
	return "", errors.New("could not load imag: file not found")
}

// hash returns a hash for the given image and options. As terminal image generation
// is based on the terminal size, the hash is based on the image hash and the options.
func hash(name string, options Options) (string, error) {
	fileHash, err := hashFile(name)
	if err != nil {
		return "", err
	}
	combined := fmt.Sprintf("%s-%s", fileHash, options.String())
	return fmt.Sprintf("%x", md5.Sum([]byte(combined))), nil
}

// getCache returns the cached data for the given hash.
func getCache(hash string) (interface{}, error) {
	path := filepath.Join("./cache", hash)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 1 {
		return base64.StdEncoding.DecodeString(lines[0])
	}

	return lo.Map(lines, func(item string, i int) string {
		res, _ := base64.StdEncoding.DecodeString(item)
		return string(res)
	}), nil
}

// setCache stores the given data in the cache directory.
func setCache(hash string, data interface{}) error {
	var lines []string

	switch d := data.(type) {
	case string:
		lines = []string{
			base64.StdEncoding.EncodeToString([]byte(d)),
		}
	case []string:
		lines = lo.Map(d, func(item string, i int) string {
			return base64.StdEncoding.EncodeToString([]byte(item))
		})
	}

	path := filepath.Join("./cache", hash)
	return ioutil.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}
