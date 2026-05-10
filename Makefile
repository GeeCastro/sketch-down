.PHONY: dev build test clean

dev:
	go run ./cmd/server/

build:
	go build -o bin/sketch-down ./cmd/server/

test:
	go test ./...

clean:
	rm -rf bin/ data/
