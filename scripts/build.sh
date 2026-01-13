#!/usr/bin/env bash

set -e

BIN=cmd/dictionary

build() {
  GOOS=$1 GOARCH=$2 OUT=$3
  echo "→ Building $OUT"
  GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
    go build -o dist/$OUT/$4 ./$BIN
}

build linux amd64 linux-amd64 dictionary
build linux arm64 linux-arm64 dictionary
build darwin amd64 darwin-amd64 dictionary
build darwin arm64 darwin-arm64 dictionary
build windows amd64 windows-amd64 dictionary.exe
