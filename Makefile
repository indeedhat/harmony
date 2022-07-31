.PHONY: build
build:
	CGO_ENABLED=0 go build -o . ./...

.PHONY: run
run: build
	./harmony-hid -v

.PHONY: run-client
run-client: build
	./harmony-client
