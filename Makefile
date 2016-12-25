NAME := godu
#VERSION := $(shell git describe --tags --abbrev=0)


## Setup
setup:
	go get github.com/cloudfoundry/bytefmt
	go get github.com/urfave/cli
	go get github.com/fatih/color
	go get github.com/Songmu/make2help/cmd/make2help
	go get gopkg.in/vmihailenco/msgpack.v2

## Run tests
test: setup
	go test ./...

## Lint
lint: setup
	golint ./...

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup test fmt help