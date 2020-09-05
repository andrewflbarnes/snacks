.SILENT:
.DEFAULT_GOAL = build
VERSION := $(shell git describe --always --long --dirty)

.PHONY: help
help:
	@echo "help:      display this help text"
	@echo "build:     build all binaries specified in cmd"

.PHONY: build
build:
	for i in cmd/*; do \
	  echo $(call goexec, build, ./$$i); \
	  $(call goexec, build, ./$$i); \
	done

.PHONY: install
install:
	@for i in cmd/*; do \
	  $(call goexec, install, ./$$i); \
	done

define goexec
	go $(strip $1) \
	  -i \
	  -v \
	  -ldflags="-X main.version=$(VERSION)" $(strip $2)
endef