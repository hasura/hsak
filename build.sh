#!/bin/bash
GOOS=darwin GOARCH=amd64 go build -o build/hsak-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o build/hsak-darwin-arm64
GOOS=linux GOARCH=amd64 go build -o build/hsak-linux-amd64
GOOS=linux GOARCH=arm64 go build -o build/hsak-linux-arm64
GOOS=windows GOARCH=amd64 go build -o build/hsak-windows-amd64.exe
