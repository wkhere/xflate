build:
	go test .
	go build

install: build
	go install -ldflags=-s

.PHONY: build install
