go:
	go vet
	go build

install: go
	go install


.PHONY: go install
