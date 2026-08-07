.DEFAULT_GOAL := check

.PHONY: fmt vet build test check race

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	go build ./...

test:
	go test ./...

race:
	go run -race ./problems/data_race

check: fmt vet build test
