.SILENT:
.DEFAULT_GOAL = build
VERSION := $(shell git describe --always --long --dirty)
TARGETS := $(wildcard cmd/*)

.PHONY: help
help:
	@echo "help:      display this help text"
	@echo "build:     (default) build all binaries specified in cmd (go build)"
	@echo "install:   install all binaries specified in cmd (go install)"

.PHONY: build
build:
	@for i in $(TARGETS); do \
	  echo $(call goexec, build, ./$$i); \
	  $(call goexec, build, ./$$i); \
	done

.PHONY: install
install:
	@for i in $(TARGETS); do \
	  echo $(call goexec, install, ./$$i); \
	  $(call goexec, install, ./$$i); \
	done

define goexec
	go $(strip $1) \
	  -v \
	  -ldflags="-X main.version=$(VERSION)" $(strip $2)
endef
