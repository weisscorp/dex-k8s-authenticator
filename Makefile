GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=dex-k8s-authenticator
REGISTRY=weisscorp
TAG=latest

all: build 

.PHONY: build
build:
	@echo "Building $(GOFILES) to ./bin"
	GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

.PHONY: container
container:
	@echo "Building container image for linux/amd64"
	docker buildx build --platform linux/amd64 -t ${GONAME}:${TAG} -t ${REGISTRY}/${GONAME}:${TAG} --load .

.PHONY: container-no-cache
container-no-cache:
	@echo "Building container image for linux/amd64 (no cache)"
	docker buildx build --platform linux/amd64 --no-cache -t ${GONAME}:${TAG} -t ${REGISTRY}/${GONAME}:${TAG} --load .
.PHONY: clean
clean:
	@echo "Cleaning"
	GOBIN=$(GOBIN) go clean
	rm -rf ./bin

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix: lint
	golangci-lint run --fix
