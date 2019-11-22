PROJECT    := github.com/cfhamlet/os-go-docid
BINNAME    ?= go-docid
BINDIR     := $(CURDIR)/bin

GOPATH      = $(shell go env GOPATH)
GOIMPORTS   = $(GOPATH)/bin/goimports
INSTALLPATH = $(GOPATH)/bin

PKG        := ./...
TESTS      := .
LDFLAGS    :=
GOFLAGS    :=
TESTFLAGS  :=
SRC        := $(shell find . -type f -name '*.go' -print)

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= $(GIT_TAG)

ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X $(PROJECT)/main.VERSION=$(BINARY_VERSION)
endif

.PHONY: test
test: build
test: TESTFLAGS += -race -v
test: test-unit
test: test-coverage

.PHONY: test-unit
test-unit:
	@echo
	@echo  "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

.PHONY: test-coverage
test-coverage:
	@echo
	@echo  "==> Running unit tests with coverage <=="
	@ ./script/coverage.sh
        

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME)


.PHONY: format
format:
	GO111MODULE=on go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w 


$(GOIMPORTS):
	(cd /; GO111MODULE=on go get -u golang.org/x/tools/cmd/goimports)

.PHONY: install
install:
	GO111MODULE=on go build -i $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(INSTALLPATH)/$(BINNAME)

.PHONY: clean
clean:
	@rm -rf $(BINDIR)

.PHONY: info
info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Tag:           ${GIT_TAG}"
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"
