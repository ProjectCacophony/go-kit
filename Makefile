.PHONY: lint

SOURCE_FOLDERS := $(shell go list -f {{.Dir}} ./...)

lint:
	golangci-lint run ./...
