.PHONY: default bench build test regen-testdata

include mks/Environment.mk

default: build
build: | test bench
test:
	go test ./...
bench:
	go test -bench=. ./...

regen-testdata:
	docker run --rm \
		-v $(shell pwd)/etf/testdata:/testdata \
		-w /testdata $(ELIXIR_IMAGE) \
		elixir /testdata/generate.exs
