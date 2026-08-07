.DEFAULT_GOAL := check

.PHONY: fmt vet build test bench check clean

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	go build ./...

test:
	go test ./...

bench:
	go test -bench=. -benchmem ./05_pointers/pointer_perf_bench ./07_generics/generics_perf_bench ./13_testing/benchmark ./14_reflect_unsafe/reflect_filter_bench

race:
	go test -tags racedemo -race ./13_testing/race_detector

check: fmt vet build test

clean:
	go clean -cache -testcache
