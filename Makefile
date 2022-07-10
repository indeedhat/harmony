.PHONY: build
build:
	go build -o . ./...

.PHONY: run
run: build
	./kvm
