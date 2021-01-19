# workdir info
PACKAGE=go-libs
PREFIX=$(shell pwd)

# which golint
GOLINT=$(shell which golangci-lint || echo '')

# build args
BUILD_ARGS :=
EXTRA_BUILD_ARGS=

export GOCACHE=

.PONY: lint test
default: lint test

lint:
	@echo "+ $@"
	@$(if $(GOLINT), , \
		$(error Please install golint: `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`))
	golangci-lint run ./...

test:
	@echo "+ test"
	go test -cover $(EXTRA_BUILD_ARGS) ./...


