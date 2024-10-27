.PHONY: build clean run

# Variables
BINARY_NAME=server
SRC_DIR=cmd
GOFILES=$(shell find . -name '*.go')

# Build the binary
build:
	@echo "Building binary..."
	go build -o bin/$(BINARY_NAME) $(SRC_DIR)/main.go

# Run the application
run: build
	@echo "Running the server..."
	./bin/$(BINARY_NAME)

# Clean up the binary
clean:
	@echo "Cleaning up..."
	rm -f bin/$(BINARY_NAME)
