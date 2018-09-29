#!/usr/bin/env bash
WASM_FILES="$(go env GOROOT)/misc/wasm"
GOOS=js GOARCH=wasm go build -o client.wasm ./client.go
cp ${WASM_FILES}/wasm_exec.js ./
go run ./server.go