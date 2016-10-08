PROJECT_NAME := vcenter-metrics
package = github.com/lenfree/$(PROJECT_NAME)

all: release

.PHONY: start
start: install
	vcenter-metrics

.PHONY: release
release: install
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/$(PROJECT_NAME)-linux-amd64 $(package)
	GOOS=darwin GOARCH=amd64 go build -o release/$(PROJECT_NAME)-darwin-amd64 $(package)

.PHONY: install
install:
	go get -v
	go fmt ./...
	go vet ./...
	go build
	go build ./...
	go install -a -tags netgo net

.PHONY: build
build: image
	docker build -t $(DOCKER_IMAGE_NAME) .

.PHONY: image
image:
ifndef DOCKER_IMAGE_NAME
	$(error DOCKER_IMAGE_NAME is not set)
endif
