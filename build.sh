#!/bin/bash

EXT=""

if [[ $GOOS == "windows" ]]; then
    EXT=".exe"
fi

go build -o project_gonzo$EXT ./cmd/game
go build -o project_gonzo_ssh$EXT ./cmd/game_ssh/
go build -o project_gonzo_browser$EXT ./cmd/game_browser