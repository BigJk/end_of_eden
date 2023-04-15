package animation

import (
	"github.com/BigJk/project_gonzo/util"
	"github.com/samber/lo"
	"math/rand"
)

var jitter = []string{"_", "#", "$", " "}

func JitterText(text string, progress float64, min int, max int) string {
	for i := 0; i < min+util.Max(0, int(float64(max)*(1-progress))); i++ {
		text = util.InsertString(text, lo.Shuffle(jitter)[0], rand.Intn(len(text)))
	}
	return text
}
