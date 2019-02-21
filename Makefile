BINARY_NAME=foodlebug.exe

build:
	go build -o $(BINARY_NAME) cmd/app.go
clean:
	go clean
	rm -f $(BINARY_NAME)
