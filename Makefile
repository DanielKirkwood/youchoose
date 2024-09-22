## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## go-lint: runs golangci-lint
.PHONY: go-lint
go-lint:
	golangci-lint run --timeout 5m

## go-format: runa go format
.PHONY: go-format
go-format:
	go fmt ./...

## go-build: run go build
.PHONY: go-build
go-build:
	go build -ldflags $(LDFLAGS) -o bin/app

## test: run tests, output by package, print coverage
.PHONY: test
test:
	@$(MAKE) go-test-by-pkg
	@$(MAKE) go-test-print-coverage

## go-test-by-pkg: run tests, output by package
.PHONY: go-test-by-pkg
go-test-by-pkg:
	gotestsum --format pkgname-and-test-fails --format-hide-empty-pkg --jsonfile /tmp/test.log -- -shuffle=on -race -cover -count=1 -coverprofile=/tmp/coverage.out ./...

## go-test-print-coverage: run tests, output by package
.PHONY: go-test-print-coverage
go-test-print-coverage:
	@printf "coverage "
	@go tool cover -func=/tmp/coverage.out | tail -n 1 | awk '{$$1=$$1;print}'

## init: run make modules and tidy
.PHONY: init
init:
	@$(MAKE) modules
	@$(MAKE) tidy

## modules: cache go modules (locally into .pkg)
.PHONY: modules
modules:
	go mod download

## tidy: tidy our go.mod file
.PHONY: tidy
tidy:
	go mod tidy -v

# only evaluated if required by a recipe
# http://make.mad-scientist.net/deferred-simple-variable-expansion/

# go module name (as in go.mod)
GO_MODULE_NAME = $(eval GO_MODULE_NAME := $$(shell \
	(mkdir -p tmp 2> /dev/null && cat tmp/.modulename 2> /dev/null) \
	|| (gsdev modulename 2> /dev/null | tee tmp/.modulename) || echo "unknown" \
))$(GO_MODULE_NAME)

# https://medium.com/the-go-journey/adding-version-information-to-go-binaries-e1b79878f6f2
ARG_COMMIT = $(eval ARG_COMMIT := $$(shell \
	(git rev-list -1 HEAD 2> /dev/null) \
	|| (echo "unknown") \
))$(ARG_COMMIT)

ARG_BUILD_DATE = $(eval ARG_BUILD_DATE := $$(shell \
	(date -Is 2> /dev/null || date 2> /dev/null || echo "unknown") \
))$(ARG_BUILD_DATE)

# https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
LDFLAGS = $(eval LDFLAGS := "\
-X '$(GO_MODULE_NAME)/internal/config.ModuleName=$(GO_MODULE_NAME)'\
-X '$(GO_MODULE_NAME)/internal/config.Commit=$(ARG_COMMIT)'\
-X '$(GO_MODULE_NAME)/internal/config.BuildDate=$(ARG_BUILD_DATE)'\
")$(LDFLAGS)
