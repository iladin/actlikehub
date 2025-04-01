# Variables
BINARY_NAME=actlikehub
SOURCE_FILE=main.go
WORKFLOW=.github/workflows/go.yml

# Build the Go program
${BINARY_NAME}: ${SOURCE_FILE}
	go build -o $@ $<

# Run the Go program with an example input file
run: ${BINARY_NAME} ${WORKFLOW}
	./$^

# Clean up the built binary
clean:
	rm -f $(BINARY_NAME)

# Run the Go program with custom input
run-custom: build
	@echo "Usage: make run-custom FILE=<path_to_yaml>"
	./$(BINARY_NAME) $(FILE)

# Default target
default: build