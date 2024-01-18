package main

import "flag"

const (
	TypeAPI         = "api"
	TypeGameContent = "game_content"
)

func main() {
	t := flag.String("type", "api", "api,game_content")
	flag.Parse()

	switch *t {
	case TypeAPI:
		buildAPIDocs()
	case TypeGameContent:
		buildGameContentDocs()
	default:
		panic("unknown type")
	}
}
