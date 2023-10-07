#!/bin/bash

for arch in "amd64" "386" "arm64" "arm"; do
    GOOS=linux GOARCH=$arch go build -ldflags "-s -w" -o ./dist/pho-$arch
    echo "Generated pho-$arch"
done
