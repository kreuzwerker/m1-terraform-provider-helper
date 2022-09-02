# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# Go variables
GO      ?= go
GOOS    ?= $(shell $(GO) env GOOS)
GOARCH  ?= $(shell $(GO) env GOARCH)
OS = $(shell uname | tr A-Z a-z)
export PATH := $(abspath bin/):${PATH}

# Project variables
COMPONENT = m1-terraform-provider-helper

# Build variables
BUILD_DIR ?= dist
VERSION = $(shell git describe --abbrev=7 --always | cut -d"-" -f 1,3 | sed -e 's/^v//g' || echo ${CI_COMMIT_REF_SLUG} )
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
LDFLAGS += -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}
export CGO_ENABLED ?= 1
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

# Dependency versions
GOTESTSUM_VERSION = 1.8.1
GOLANGCI_VERSION = 1.49.0
GITCHGLOG_VERSION = 0.15.1
SVU_VERSION = 1.9.0

GOLANG_VERSION = 1.19

# Add the ability to override some variables
# Use with care
-include override.mk

.PHONY: clear
clear: ## Clear the working area and the project
	rm -rf bin data dist

.PHONY: goversion
goversion:
ifneq (${IGNORE_GOLANG_VERSION_REQ}, 1)
	@printf "${GOLANG_VERSION}\n$$(go version | awk '{sub(/^go/, "", $$3);print $$3}')" | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -g | head -1 | grep -q -E "^${GOLANG_VERSION}$$" || (printf "Required Go version is ${GOLANG_VERSION}\nInstalled: `go version`" && exit 1)
endif

.PHONY: appversion
appversion: ## Print app-version to stdout
	@echo "${VERSION}"

.PHONY: build-%
build-%: goversion
ifeq (${VERBOSE}, 1)
	go env
endif

	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/$* ./cmd/$*

.PHONY: build
build: goversion ## Build all binaries
ifeq (${VERBOSE}, 1)
	go env
endif

	@mkdir -p ${BUILD_DIR}
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/ ./cmd/...

.PHONY: build-release
build-release: ## Build all binaries without debug information
	@${MAKE} LDFLAGS="-w ${LDFLAGS}" GOARGS="${GOARGS} -trimpath" BUILD_DIR="${BUILD_DIR}/release" build

.PHONY: build-debug
build-debug: ## Build all binaries with remote debugging cabilities
	@${MAKE} GOARGS="${GOARGS} -gcflags \"all=-N -l\"" BUILD_DIR="${BUILD_DIR}/debug" build

.PHONY: generate-mocks
generate-mocks: ## regenerates the mocks for the tests
	mockery --all

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${GOOS}_${GOARCH}.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}

TEST_PKGS ?= ./...
TEST_REPORT_NAME ?= results.xml
.PHONY: test
test: TEST_REPORT ?= main
test: TEST_FORMAT ?= short-verbose
test: TEST_TIMEOUT ?= 1m
test: SHELL = /bin/bash
test: bin/gotestsum ## General test cmd. DO NOT use directly, but use test-functional or test-integration
	@mkdir -p ${BUILD_DIR}/test_results/${TEST_REPORT}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/test_results/${TEST_REPORT}/${TEST_REPORT_NAME} --jsonfile ${BUILD_DIR}/test_results/${TEST_REPORT}/results.out --format ${TEST_FORMAT} -- -timeout ${TEST_TIMEOUT} -coverprofile=${BUILD_DIR}/test_results/${TEST_REPORT}/coverage.out -covermode atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

.PHONY: test-functional
test-functional: ## Run functional tests
	@${MAKE} GOARGS="${GOARGS} -failfast -race -run ^.*\$$\$$" TEST_REPORT=functional test


bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run

.PHONY: lint-ci
lint-ci: bin/golangci-lint ## Run linter with output to file with checkstyle format for sonar
	@mkdir -p ${BUILD_DIR}/lint_results
	bin/golangci-lint run --out-format checkstyle > ${BUILD_DIR}/lint_results/golangci-lint.out

bin/svu: bin/svu-${SVU_VERSION}
	@ln -sf svu-${SVU_VERSION} bin/svu
bin/svu-${SVU_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/caarlos0/svu/releases/download/v${SVU_VERSION}/svu_${SVU_VERSION}_${GOOS}_${GOARCH}.tar.gz | tar -zOxf - svu > ./bin/svu-${SVU_VERSION} && chmod +x ./bin/svu-${SVU_VERSION}

bin/git-chglog: bin/git-chglog-${GITCHGLOG_VERSION}
	@ln -sf git-chglog-${GITCHGLOG_VERSION} bin/git-chglog
bin/git-chglog-${GITCHGLOG_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/git-chglog/git-chglog/releases/download/v${GITCHGLOG_VERSION}/git-chglog_${GITCHGLOG_VERSION}_${GOOS}_${GOARCH}.tar.gz | tar -zOxf - git-chglog > ./bin/git-chglog-${GITCHGLOG_VERSION} && chmod +x ./bin/git-chglog-${GITCHGLOG_VERSION}

release-%: TAG_PREFIX = ""
release-%: bin/git-chglog
	@echo "Generating CHANGELOG"
	bin/git-chglog --next-tag $* -o CHANGELOG.md
ifeq (${TAG}, 1)
	@echo "Committing and tagging"
	git add CHANGELOG.md
	git commit -m 'Prepare release $*'
	git tag -m 'Release $*' ${TAG_PREFIX}$*
ifeq (${PUSH}, 1)
	git push; git push origin ${TAG_PREFIX}$*
endif
endif

	@echo "Version updated to $*!"
ifneq (${PUSH}, 1)
	@echo
	@echo "Review the changes made by this script then execute the following:"
ifneq (${TAG}, 1)
	@echo
	@echo "git add CHANGELOG.md cmd/version.go && git commit -m 'Prepare release $*' && git tag -m 'Release $*' ${TAG_PREFIX}$*"
	@echo
	@echo "Finally, push the changes:"
endif
	@echo
	@echo "git push; git push origin ${TAG_PREFIX}$*"
endif

.PHONY: patch
patch: ## Release a new patch version
	@${MAKE} replace-occurences-$(shell (svu --strip-prefix patch))
	@${MAKE} release-$(shell (svu --strip-prefix patch))

.PHONY: minor
minor: ## Release a new minor version
	@${MAKE} replace-occurences-$(shell (svu --strip-prefix minor))
	@${MAKE} release-$(shell (svu --strip-prefix minor))

.PHONY: major
major: ## Release a new major version
	@${MAKE} replace-occurences-$(shell (svu --strip-prefix major))
	@${MAKE} release-$(shell (svu --strip-prefix major))

replace-occurences-%:
	@echo "Replace occurences of old version strings..."
	sed -i '' "s/$(shell (svu --strip-prefix current))/$*/g" cmd/version.go

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)