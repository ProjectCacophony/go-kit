.PHONY: lint

SOURCE_FOLDERS := $(shell go list -f {{.Dir}} ./...)

lint:
	goimports -d $(SOURCE_FOLDERS)
	golangci-lint run --deadline=30m --enable-all ./...
