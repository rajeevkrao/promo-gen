start:
	go run main.go

build:
	go build -o build/PromoGen.exe

test:
	go test ./...

clean:
	del myapp.exe