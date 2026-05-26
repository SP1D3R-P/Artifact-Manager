

SRC_DIR  := internal
CMD_DIR  := cmd

DIRS     := $(CMD_DIR) $(SRC_DIR)
GO_FILES := $(foreach dir,$(DIRS),$(wildcard $(dir)/*.go) $(wildcard $(dir)/**/*.go))


.PHONY: all build run test docker-build clean

all: build
	@echo "Build complete. Run 'make run' to execute."

$(BUILDER_BIN_PATH): $(GO_FILES)
	@echo "Building builder..."
	go build -o $(BUILDER_BIN_PATH) ./$(CMD_DIR)

build: $(BUILDER_BIN_PATH)

test:
	@echo "Running tests..."
	go test -v -cover ./...

docker-build: $(BUILDER_BIN_PATH)
	@echo "Building Docker image..."
	docker build -t builder:latest .






