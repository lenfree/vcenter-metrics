all: start

start:
	go get -v
	go fmt ./...
	go vet ./...
	go build
	go build ./...
	go install -a -tags netgo net
	vcenter-metrics

.PHONY: start
