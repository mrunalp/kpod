GO ?= go
EPOCH_TEST_COMMIT ?= 736c36c1a11b690b105056bdd627a37c97011ffd 
PROJECT := github.com/mrunalp/kpod
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_BRANCH_CLEAN := $(shell echo $(GIT_BRANCH) | sed -e "s/[^[:alnum:]]/-/g")
KPOD_IMAGE := kpod_dev$(if $(GIT_BRANCH_CLEAN),:$(GIT_BRANCH_CLEAN))
KPOD_LINK := ${CURDIR}/vendor/src/github.com/mrunalp/kpod
KPOD_LINK_DIR := ${CURDIR}/vendor/src/github.com/mrunalp
KPOD_INSTANCE := kpod_dev
SYSTEM_GOPATH := ${GOPATH}
PREFIX ?= ${DESTDIR}/usr
BINDIR ?= ${PREFIX}/bin
LIBEXECDIR ?= ${PREFIX}/libexec
MANDIR ?= ${PREFIX}/share/man
ETCDIR ?= ${DESTDIR}/etc
export GOPATH := ${CURDIR}/vendor

all: binaries

default: help

help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'binaries' - Build kpod"
	@echo " * 'clean' - Clean artifacts"
	@echo " * 'lint' - Execute the source code linter"

lint: ${KPOD_LINK}
	@which gometalinter > /dev/null 2>/dev/null || (echo "ERROR: gometalinter not found. Consider 'make install.tools' target" && false)
	@echo "checking lint"
	@./.tool/lint


${KPOD_LINK}:
	mkdir -p ${KPOD_LINK_DIR}
	ln -sfn ${CURDIR} ${KPOD_LINK}


GO_SRC =  $(shell find . -name \*.go)

kpod: $(GO_SRC) | ${KPOD_LINK}
	$(GO) build --tags "$(BUILDTAGS)" -o $@ 

clean:
	rm -f kpod
	rm -f ${KPOD_LINK}
	find . -name \*~ -delete
	find . -name \#\* -delete

binaries: kpod

.PHONY: .gitvalidation
# When this is running in travis, it will only check the travis commit range
.gitvalidation:
	@which git-validation > /dev/null 2>/dev/null || (echo "ERROR: git-validation not found. Consider 'make install.tools' target" && false)
ifeq ($(TRAVIS),true)
	git-validation -q -run DCO,short-subject
else
	git-validation -v -run DCO,short-subject -range $(EPOCH_TEST_COMMIT)..HEAD
endif

.PHONY: install.tools

install.tools: .install.gitvalidation .install.gometalinter .install.md2man

.install.gitvalidation:
	GOPATH=${SYSTEM_GOPATH} go get github.com/vbatts/git-validation

.install.gometalinter:
	GOPATH=${SYSTEM_GOPATH} go get github.com/alecthomas/gometalinter
	GOPATH=${SYSTEM_GOPATH} gometalinter --install

.install.md2man:
	GOPATH=${SYSTEM_GOPATH} go get github.com/cpuguy83/go-md2man

.PHONY: \
	binaries \
	clean \
	lint
