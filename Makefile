GO ?= go
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


.PHONY: \
	binaries \
	clean
