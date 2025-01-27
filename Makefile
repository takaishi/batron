.PHONY: build
build:
	go build -o dist/batron ./cmd/batron

.PHONY: install
install:
	go install github.com/takaishi/batron/cmd/batron

.PHONY: test
test:
	go test -race ./...
