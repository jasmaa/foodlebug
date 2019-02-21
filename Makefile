BINARY_NAME=foodlebug.exe

start:
	go build -o $(BINARY_NAME) cmd/app.go
	./$(BINARY_NAME)
build:
	go build -o $(BINARY_NAME) cmd/app.go
clean:
	go clean
	rm -f $(BINARY_NAME)
