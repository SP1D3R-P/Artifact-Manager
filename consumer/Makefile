.PHONY: all build run test docker-build clean

CMD_DIR  := ./cmd
SRC_DIR  := internal
DIRS     := $(CMD_DIR) $(SRC_DIR)

GO_FILES := $(foreach dir,$(DIRS),$(wildcard $(dir)/*.go) $(wildcard $(dir)/**/*.go))


$(CONSUMER_BIN_PATH): $(GO_FILES)
	@echo "Building the consumer..."
	go build -o $(CONSUMER_BIN_PATH) $(CMD_DIR)/
	
build: $(CONSUMER_BIN_PATH)

test:
	@echo "Running tests..."
	go test ./...

docker-build: 
	@echo "Building Docker image for the consumer..."
	docker build -t consumer:latest .

help : 
	@echo "hello"

.DEFAULT := help