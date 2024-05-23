ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN=$(ROOT_DIR)/bin/memdb
REVISION=$(shell git describe --tags 2>/dev/null || git log --format="v0.0-%h" -n 1 || echo "v0.0-unknown")
SHELL=/bin/bash -euo pipefail

prepare: config.yml
	@true

deps:
	@go mod vendor

build:
	@go build -trimpath -ldflags "-X memdb.revision=$(REVISION)" -o $(BIN) $(ROOT_DIR)/cmd/.

.PHONY: test
test:
	@go install gotest.tools/gotestsum@latest
	@gotestsum --no-color=false --junitfile $(ROOT_DIR)/junit-report.xml -- ./... -v -race -parallel 8 -count 1 -coverpkg=./... -coverprofile=$(ROOT_DIR)/coverage.out

coverage-report: coverage.out
	@go tool cover -func=$(ROOT_DIR)/coverage.out
	@go install github.com/boumenot/gocover-cobertura@latest
	@gocover-cobertura < $(ROOT_DIR)/coverage.out > $(ROOT_DIR)/coverage.xml

#.env:
#	@cp .env.dist .env

config.yml:
	@cp config/config.yaml.dist config/config.yaml
	@cp config/config_test.yaml.dist config/config_test.yaml