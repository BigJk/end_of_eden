package clipboard

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestSet(t *testing.T) {
	Set("hello world")
	res, _ := exec.Command("pbpaste").CombinedOutput()
	assert.Equal(t, "hello world", string(res))
}
