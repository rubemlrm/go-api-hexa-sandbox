TARGET_BIN = go-api-template
TARGET_ARCH = amd64
SOURCE_MAIN = cmd/app/main.go
LDFLAGS = -s -w

all: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(TARGET_ARCH) go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN)_linux-amd64 $(SOURCE_MAIN)

build-linux-noarch:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN) $(SOURCE_MAIN)

start:
	go run $(SOURCE_MAIN)

install-dependencies:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4

.PHONY: run
mod-download:
	go mod download

generate: install-dependencies mod-download
	go generate ./...
