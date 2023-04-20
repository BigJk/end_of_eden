#!/bin/bash

go build -o project_gonzo ./cmd/game
go build -o project_gonzo_ssh ./cmd/game_ssh/
go build -o project_gonzo_browser ./cmd/game_browser