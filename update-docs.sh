#!/bin/bash

go run ./cmd/docs > ./docs/LUA_API_DOCS.md
go run ./cmd/definitions > ./assets/scripts/definitions/api.lua