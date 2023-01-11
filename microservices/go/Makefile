.DEFAULT_GOAL := help

# File Paths
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
CMD_DIR := $(ROOT_DIR)/cmd
PKG_DIR := $(ROOT_DIR)/pkg
BIN_DIR := $(ROOT_DIR)/bin
MAIN_FILE := main.go

# Selected Service File Paths
SELECTED_SERVICE_DIR := $(CMD_DIR)/$(service)
SELECTED_SERVICE_ROOT_PKG := github.com/obenkenobi/cypher-log/microservices/go/cmd/$(service)
SELECTED_SERVICE_DIR_MAIN := $(SELECTED_SERVICE_DIR)/$(MAIN_FILE)
SELECTED_TEST_LOCATION=$(testlocation)
TEST_PATH=$(ROOT_DIR)/$(SELECTED_TEST_LOCATION)

# Go Architectures
GOARCH_ARM64 := arm64
GOARCH_AMD64 := amd64

# Go Operating Systems
GOOS_LINUX := linux
GOOS_WINDOWS := windows
GOOS_DARWIN := darwin

# Go binaries
EXECUTABLE := $(service).out

# Go binaries (full paths)
EXECUTABLE_PATH := $(BIN_DIR)/$(EXECUTABLE)

define go_build_platform
	env GOOS=$(1) GOARCH=$(2) go build -o $(3) $(4)
endef

.PHONY: help
help:
	@echo "Arguments:"
	@echo ""
	@echo "service			a microservice in the cmd directory"
	@echo "testlocation  	path to the directory or file you want to run a go test on"
	@echo ""
	@echo "targets:"
	@echo ""
	@echo "(in the format <target> : <dependencies> |<required args>| <description>"
	@echo ""
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: list-microservices
list-microservices: ##|| Lists available microservices
	ls $(CMD_DIR)

.PHONY: wire
wire: ##|service| Code generation for dependency injection
	@cd $(SELECTED_SERVICE_DIR)
	go run github.com/google/wire/cmd/wire gen $(SELECTED_SERVICE_ROOT_PKG)/app

.PHONY: wire-all
wire-all: $(CMD_DIR)/* ##|| Code generation for all microservices, dependency injection
	@for microservice in $(shell ls ${CMD_DIR}); do \
  		cd $(CMD_DIR)/$${microservice}; \
		go run github.com/google/wire/cmd/wire gen github.com/obenkenobi/cypher-log/microservices/go/cmd/$${microservice}/app ; \
    done

.PHONY: build
build: wire ##|service| builds your executable in the bin directory in the format $(service).out
	$(call go_build_platform,,,$(EXECUTABLE_PATH),$(SELECTED_SERVICE_DIR))

.PHONY: build-linux-amd64
build-linux-amd64: wire ##|service| builds your executable in the bin directory in the format $(service).out for the linux/amd64 platform
	$(call go_build_platform,$(GOOS_LINUX),$(GOARCH_AMD64),$(EXECUTABLE_PATH),$(SELECTED_SERVICE_DIR))

.PHONY: build-linux-arm64
build-linux-arm64: wire ##|service| builds your executable in the bin directory in the format $(service).out for the linux/arm64 platform
	$(call go_build_platform,$(GOOS_LINUX),$(GOARCH_ARM64),$(EXECUTABLE_PATH),$(SELECTED_SERVICE_DIR))

.PHONY: build-all
build-all: wire-all $(CMD_DIR)/* ##|| for all microservices builds your executables in the bin directory in the format $(service).out
	@for microservice in $(shell ls ${CMD_DIR}); do \
  		go build -o $(BIN_DIR)/$${microservice}.out $(CMD_DIR)/$${microservice} ; \
    done

.PHONY: run
run: build ##|service| runs your go executable
	exec $(EXECUTABLE_PATH)

.PHONY: docker-linux-amd64 ##|service| builds a docker image for the linux/amd64 platform
docker-linux-amd64: build-linux-amd64
	docker build -t $(service)_$(GOOS_LINUX)_$(GOARCH_AMD64):latest --build-arg EXECUTABLE=$(EXECUTABLE) $(ROOT_DIR)/

.PHONY: docker-linux-arm64
docker-linux-arm64: build-linux-arm64 ##|service| builds a docker image for the linux/arm64 platform
	docker build build -t $(service)_$(GOOS_LINUX)_$(GOARCH_ARM64):latest --build-arg EXECUTABLE=$(EXECUTABLE) $(ROOT_DIR)/

.PHONY: test
test: wire-all ##|testlocation| runs go tests
ifeq ($(testlocation),)
	go test $(ROOT_DIR)/...
else ifeq ($(suffix $(TEST_PATH)), .go)
	go test $(TEST_PATH)
else
	go test $(TEST_PATH)/...
endif

clean:
	rm -rf $(BIN_DIR)/*