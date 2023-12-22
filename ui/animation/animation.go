package animation

import (
	"github.com/BigJk/end_of_eden/ui"
	"github.com/samber/lo"
	"math/rand"
)

var jitter = []string{"_", "#", "$", " "}

func JitterText(text string, progress float64, min int, max int) string {
	for i := 0; i < min+ui.Max(0, int(float64(max)*(1-progress))); i++ {
		text = ui.InsertString(text, lo.Shuffle(jitter)[0], rand.Intn(len(text)))
	}
	return text
}
