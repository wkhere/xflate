go:
	go test
	go build

install: go
	go install


.PHONY: go install
