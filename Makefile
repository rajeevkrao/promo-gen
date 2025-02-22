.PHONY: run build test

start:
	go run .

build:
	go build -o build/PromoGen.exe

test:
	go test ./...