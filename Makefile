.PHONY: full build test test-go lint lint-go fix fix-go clean

SHELL=/bin/bash -o pipefail
$(shell git config core.hooksPath ops/git-hooks)
GO_PATH := $(shell go env GOPATH 2> /dev/null)
PATH := /usr/local/bin:$(GO_PATH)/bin:$(PATH)

full: clean lint test build

## Build the project
build:

## Test the project
test: test-go

test-go:
	@mkdir -p var/coverage/go/
	@go install github.com/boumenot/gocover-cobertura@latest
	go test -race -cover -coverprofile var/coverage/go/profile.txt ./...
	@go tool cover -func var/coverage/go/profile.txt | awk '/^total/{print $$1 " " $$3}'
	@go tool cover -html var/coverage/go/profile.txt -o var/coverage/go/coverage.html
	@gocover-cobertura < var/coverage/go/profile.txt > var/coverage/go/cobertura-coverage.xml

## Lint the project
lint: lint-go

lint-go:
	@go install golang.org/x/tools/cmd/goimports@latest
	go get -d ./...
	go mod tidy
	gofmt -s -w .
	go vet ./...
	goimports -w .

## Fix the project
fix: fix-go

fix-go:
	go mod tidy
	gofmt -s -w .
	goimports -w .

## Clean the project
clean:
	git clean -Xdff --exclude="!.env*local"
