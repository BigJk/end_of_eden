#!/bin/bash

echo "Updating docs..."
go run ./cmd/internal/docs > ./docs/LUA_API_DOCS.md
echo "Done!"

echo "Updating game docs..."
go run ./cmd/internal/docs -type game_content > ./docs/GAME_CONTENT_DOCS.md
echo "Done!"

echo "Updating definitions..."
go run ./cmd/internal/definitions > ./assets/scripts/definitions/api.lua
echo "Done!"
