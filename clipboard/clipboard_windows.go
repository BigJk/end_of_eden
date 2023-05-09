package clipboard

import clip "golang.design/x/clipboard"

func Init() {
	if err := clip.Init(); err != nil {
		panic(err)
	}
}

func Set(data string) {
	clip.WriteAll(data)
}
