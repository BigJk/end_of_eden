package faces

import (
	"fmt"
	"github.com/BigJk/end_of_eden/fs"
	"github.com/samber/lo"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
)

// FaceGenerator is a generator for faces.
type FaceGenerator struct {
	data map[int][][]string
}

// Gen generates a face with the given id.
func (gen *FaceGenerator) Gen(id int) string {
	var face []string

	t := gen.data[id]
	for i := 0; i < 7; i++ {
		if len(t[i]) == 0 {
			continue
		}
		face = append(face, t[i][rand.Intn(len(t[i]))])
	}

	minSpace := lo.Min(lo.Map(face, func(line string, _ int) int {
		count := 0
		for _, v := range line {
			if v == ' ' {
				count++
			} else {
				break
			}
		}

		return count
	}))

	if minSpace > 0 {
		for i := range face {
			face[i] = face[i][minSpace:]
		}
	}

	return strings.Join(face, "\n")
}

// GenRand generates a random face.
func (gen *FaceGenerator) GenRand() string {
	if gen == nil || gen.data == nil || len(gen.data) == 0 {
		return ""
	}
	return gen.Gen(lo.Shuffle(lo.Keys(gen.data))[0])
}

// New creates a new FaceGenerator.
func New(dataFolder string) (*FaceGenerator, error) {
	gen := &FaceGenerator{
		data: map[int][][]string{},
	}
	for i := 0; i < 7; i++ {
		bytes, err := fs.ReadFile(filepath.Join(dataFolder, fmt.Sprintf("/Face%d.txt", i)))
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(bytes), "\r\n")
		for j := range lines {
			split := strings.SplitN(lines[j], ".", 2)
			id, _ := strconv.ParseInt(split[0], 10, 64)

			if len(split) != 2 {
				continue
			}

			if _, ok := gen.data[int(id)]; !ok {
				gen.data[int(id)] = [][]string{{}, {}, {}, {}, {}, {}, {}}
			}

			gen.data[int(id)][i] = append(gen.data[int(id)][i], split[1])
		}
	}
	return gen, nil
}

// Global is the global generator.
var Global *FaceGenerator

// InitGlobal initializes the global generator.
func InitGlobal(dataFolder string) error {
	gen, err := New(dataFolder)
	if err != nil {
		return err
	}
	Global = gen
	return nil
}
