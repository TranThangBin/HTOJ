#!/usr/bin/env sh

./tailwindcss-linux-x64 -i ./src/input.css -o ./public/style.css
go tool templ generate ./templates
go mod tidy
go build ./cmd/HTOJ