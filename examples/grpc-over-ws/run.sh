#!/usr/bin/env bash
WASM_FILES="$(go env GOROOT)/misc/wasm"
GOOS=js GOARCH=wasm go build -o client.wasm ./client.go || exit 1
cp ${WASM_FILES}/wasm_exec.js ./
go run ./server.go