#!/bin/bash

EXT=""

if [[ $GOOS == "windows" ]]; then
    EXT=".exe"
fi

go build -o end_of_eden$EXT ./cmd/game/
go build -o end_of_eden_ssh$EXT ./cmd/game_ssh/
go build -o end_of_eden_browser$EXT ./cmd/game_browser/