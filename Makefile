.DEFAULT_GOAL := default

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))


.PHONY: all
all: help
	@:

IMAGE ?= quay.io/kmamgain/event-listener:latest

export DOCKER_CLI_EXPERIMENTAL=enabled

.PHONY: build # Build the container image
build:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
		--output "type=docker,push=false" \
		--tag $(IMAGE) \
		.

.PHONY: publish # Push the image to the remote registry
publish:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--output "type=image,push=true" \
		--tag $(IMAGE) \
		.

.PHONY: build-go
build-go:
	go build .