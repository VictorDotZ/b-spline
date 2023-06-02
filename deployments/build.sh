#!/usr/bin/bash

export PATH=$PATH:~/tmp/go/bin

go build -o 1d.out ./cmd/1d/main.go
go build -o 2d.out ./cmd/2d/main.go
