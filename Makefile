VERSION := $(shell cat version)
PACKAGE_NAME := $(shell go list .)
BINDIR := bin

.PHONY: setup
setup:
	@go install github.com/mitchellh/gox

.PHONY: build
build:
	-@make setup
	-@gox -output "bin/{{.Dir}}-$(VERSION)-{{.OS}}-{{.Arch}}"
