#!/bin/sh

# The Environments
GO111MODULE=on
GOSUMDB=sum.golang.google.cn
GOPROXY=https://goproxy.cn,direct

# Version Information
COMMIT=$(git rev-parse HEAD 2>/dev/null)
VERSION=$(git describe --tags 2>/dev/null)
BUILD_DATE=$(date +"%s")

# Build App
go build -ldflags "-w -X github.com/xgfone/gover.Commit=$COMMIT -X github.com/xgfone/gover.Version=$VERSION -X github.com/xgfone/gover.BuildTime=$BUILD_DATE" github.com/xgfone/bttools/...
