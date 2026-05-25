
CURR_DIR := $(shell pwd ) 

BIN := bin
BIN_PATH :=$(strip $(CURR_DIR))/$(strip $(BIN))


BUILDER := builder
CONSUMER := consumer

MODULES := $(BUILDER) $(CONSUMER)

# Tests 
TEST_DATAS := testData/*
TEST_DATA_DIRS :=  $(shell find $(TEST_DATAS) -maxdepth 0 -type d )
TESTS := $(sort $(foreach d,$(TEST_DATA_DIRS),$(if $(wildcard $(d)/config.json),$(d),)))

.PHONY : test all 

help : 
	@echo "Available targets:"
	@echo "  docker-build  - Build Docker images for all modules"
	@echo "  k8-start      - Start Kubernetes services after building Docker images"
	

docker-build: 
	@echo "Building Docker images for all modules..."
	@for module in $(MODULES); do \
		$(MAKE) -C $$module docker-build; \
	done

build:


k8-start: docker-build
	cd k8s && make start