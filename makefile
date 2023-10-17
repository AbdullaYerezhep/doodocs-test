# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARY_NAME = test

# Directories
BUILD_DIR = build
CMD_DIR = cmd
INTERNAL_DIR = internal

# Targets
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go

run:
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@$(GOTEST) -v ./$(INTERNAL_DIR)/handler

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

.PHONY: build run test clean
