# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

.DEFAULT_GOAL := build

BINARY_NAME:=lh-cdc-singularity.sif
BINARIES:=${BINARY_NAME}

# Ensure that make is run with bash shell as some syntax below is bash-specific.
SHELL:=/usr/bin/env bash

# Set GOPRIVATE if it is not already set.
ifndef GOPRIVATE
export GOPRIVATE="github.com/hpe-hcss"
endif

ifndef PARENT_BRANCH_NAME
PARENT_BRANCH_NAME:=main
PARENT_BRANCH_SHA:=$(shell git merge-base main HEAD)
endif

# Directories
export GOBIN=${PWD}/bin

# The location of the go mock generator. It must be exported so it will work
# with go generate.
export GO_MOCKGEN := ${GOBIN}/mockgen

# The location of any generated golang mocks.
GO_MOCK_DIR := internal/pkg/mocks

# The location to which test reports are written.
TESTREPORT_DIR := test-reports

# The location to which coverage reports are written. Changing this value will
# likely cause CircleCI failures.
COVERAGE_DIR := coverage/go

# The markers have to live in the directory they are marking, otherwise
# they may not get carried through operations such as a Dockerfile COPY.
MAKE_MARKER := .make_marker
GO_MOCK_MARKER:=${GO_MOCK_DIR}/${MAKE_MARKER}

# download module dependencies
.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download

# Install command dependencies.
.PHONY: tools
tools: download
	@echo Installing tools from tools.go
	@go mod download github.com/golang/mock # Need to check if that needs to be done in the project 
	@awk -F'"' '/_/{print $$2}' < tools.go | xargs -tI % go install %

# Generate the go mocks. Note that this rule does not have the dependencies
# correct for mocks build within this git repository. Changes to mocked
# interfaces here will not cause the mock to be rebuilt.
${GO_MOCK_MARKER}: $(wildcard ${GO_MOCK_DIR}/*/generate.go) | tools
	@mkdir -p ${GO_MOCK_DIR}
	@cd ${GO_MOCK_DIR} && go generate ./...
	@touch ${GO_MOCK_MARKER}

# mocks is a target to manually rebuild the mocks without having to remember the
# name of the marker. Because it is a phoney target, it always says it needs to
# be rebuilt, so other make targets should depend on the markers instead.
.PHONY: mocks
mocks: ${GO_MOCK_MARKER}

.PHONY: generate
generate:  mocks

.PHONY: build
build: ${BINARIES}

${BINARY_NAME}: generate
	echo  building singularity plugin
	singularity plugin compile .
	ls /root/project

.PHONY: package
package: build 
	pwd
	ls /root/project
	rm -rf /root/project/tarballs 
	mkdir -p /root/project/tarballs 
	cd /go/ && tar jcf /root/project/tarballs/singularity-${SINGULARITY_VERSION}.tbz singularity
	cd /root/project/tarballs && cp  /root/project/project.sif ./lh-cdc-singularity.sif 
	cd /root/project/tarballs && tar jcf  ./lh-cdc-singularity-${CIRCLE_TAG}.tbz ./singularity-${SINGULARITY_VERSION}.tbz   ./lh-cdc-singularity.sif 
	cd /root/project/tarballs && rm -f ./lh-cdc-singularity.sif ./singularity-${SINGULARITY_VERSION}.tbz


# Run go test against all go source files
.PHONY: test
test: ${GO_MOCK_MARKER} build
	go mod tidy
	gotestsum --format short-verbose -- ./... --race

# compiling a singularity module does not generate a go file so in order for the linter to function we are building the test. 
.PHONY: lint 
lint: ${GO_MOCK_MARKER}
	go get ./... 
	golangci-lint run

.PHONY: lint-new 
lint-new: ${GO_MOCK_MARKER}
	echo "Showing lint issues added after ${PARENT_BRANCH_NAME} (${PARENT_BRANCH_SHA})"
	go get ./...
	golangci-lint run --new-from-rev ${PARENT_BRANCH_SHA}

.PHONY: coverage
coverage: ${GO_MOCK_MARKER} build
	mkdir -p ${COVERAGE_DIR}/html
	echo "generating coverage data"
	go mod tidy
	go test -short -race -coverprofile=${COVERAGE_DIR}/coverage.tmp ./...
	grep -v -e mock -e generated < ${COVERAGE_DIR}/coverage.tmp > ${COVERAGE_DIR}/coverage.out
	go tool cover -html=${COVERAGE_DIR}/coverage.out -o ${COVERAGE_DIR}/html/main.html
	echo "Generated ${COVERAGE_DIR}/html/main.html"

.PHONY: run
run:
	./${BINARY_NAME}

.PHONY: clean
clean:
	go clean
	rm -rf ${BINARY_NAME} ${GOBIN} ${COVERAGE_DIR}
	find . -name "generated*" | xargs rm -f
	find . -name ${MAKE_MARKER} | xargs rm -f

