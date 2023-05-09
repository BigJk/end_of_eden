package clipboard

import "github.com/atotto/clipboard"

func Init() {
	if err := clipboard.Init(); err != nil {
		panic(err)
	}
}

func Set(data string) {
	clipboard.WriteAll(data)
}
