TARGET_BIN = go-api-template
TARGET_ARCH = amd64
SOURCE_MAIN = cmd/app/main.go
LDFLAGS = -s -w

export GOOSE_DBSTRING=postgresql://demo:demo@127.0.0.1:5432/demo
export GOOSE_DRIVER=postgres

all: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(TARGET_ARCH) go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN)_linux-amd64 $(SOURCE_MAIN)

build-linux-noarch:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN) $(SOURCE_MAIN)

start:
	go run $(SOURCE_MAIN)

install-dependencies:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest


.PHONY: run
mod-download:
	go mod download

generate: install-dependencies mod-download
	oapi-codegen -generate types -o internal/user/ports/openapi_types.gen.go -package ports spec/user.yaml
	oapi-codegen -generate gin-server -o internal/user/ports/openapi_api.gen.go -package ports spec/user.yaml

generate-mocks:
	@mockery --output user/mocks --dir user --all

.PHONY: build
build: ## Build app
	go build -o bin/app cmd/app/main.go


migrate: ## run database migrations
	goose -dir migrations up


migrate-rollback: ## run database migrations
	goose -dir migrations down
