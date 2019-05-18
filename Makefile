.PHONY: lint

SOURCE_FOLDERS := $(shell go list -f {{.Dir}} ./...)

lint:
	golangci-lint run --deadline=30m --disable-all \
	--enable=govet \
	--enable=staticcheck \
	--enable=unused \
	--enable=gosimple \
	--enable=structcheck \
	--enable=varcheck \
	--enable=ineffassign \
	--enable=deadcode \
	--enable=golint \
	--enable=unconvert \
	--enable=goimports \
	--enable=maligned \
	--enable=unparam \
	--enable=prealloc \
	--enable=scopelint \
	./...
