# Golang Env Variables
GO111MODULE=on
GOSUMDB=sum.golang.google.cn
GOPROXY=https://goproxy.cn,direct

# Build Version Information
COMMIT=$(shell git rev-parse HEAD 2>/dev/null)
VERSION=$(shell git describe --tags 2>/dev/null)
BUILD_DATE=$(shell date +"%s")

BUILD_FLAGS_DATE=-X github.com/xgfone/gover.BuildTime=$(BUILD_DATE)
BUILD_FLAGS_COMMIT=-X github.com/xgfone/gover.Commit=$(COMMIT)
BUILD_FLAGS_VERSION=-X github.com/xgfone/gover.Version=$(VERSION)
BUILD_FLAGS_X=$(BUILD_FLAGS_DATE) $(BUILD_FLAGS_COMMIT) $(BUILD_FLAGS_VERSION)

.PHONY: all install build download
all: build

install: download
	go install -ldflags "-w $(BUILD_FLAGS_X)"

build: download
	go build -ldflags "-w $(BUILD_FLAGS_X)"

download:
	go mod download
