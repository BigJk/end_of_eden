#!/bin/bash

echo "Updating docs..."
go run ./cmd/internal/docs > ./docs/LUA_API_DOCS.md
echo "Done!"

echo "Updating definitions..."
go run ./cmd/internal/definitions > ./assets/scripts/definitions/api.lua
echo "Done!"
