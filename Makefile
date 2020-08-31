.SILENT:
.DEFAULT_GOAL = build

.PHONY: help
help:
	@echo "help:      display this help text"
	@echo "build:     build all binaries specified in cmd"

.PHONY: build
build:
	@for i in cmd/*; do go build ./$$i; done

.PHONY: install
install:
	@for i in cmd/*; do go install ./$$i; done