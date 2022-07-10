.PHONY: build
build:
	go build -o . ./...

.PHONY: run
run: build
	./harmony-server

.PHONY: run-client
run-client: build
	./harmony-client
