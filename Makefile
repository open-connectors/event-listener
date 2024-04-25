.DEFAULT_GOAL := default

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))


.PHONY: all
all: help
	@:

IMAGE ?= quay.io/kmamgain/cdevent:latest

export DOCKER_CLI_EXPERIMENTAL=enabled

.PHONY: build # Build the container image
build:
	podman build --platform linux/amd64,linux/arm64 -t quay.io/kmamgain/cdevent:latest  .

.PHONY: publish # Push the image to the remote registry
publish:
	podman push $(IMAGE)

.PHONY: build-go
build-go:
	go build .