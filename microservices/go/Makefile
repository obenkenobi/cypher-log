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

# Go binary extensions
EXTENSION_BIN_LINUX := linux
EXTENSION_BIN_DARWIN := darwin
EXTENSION_BIN_WINDOWS := exe

# Go binaries (full paths)
BIN_LINUX_AMD64 := $(BIN_DIR)/$(service).$(GOARCH_AMD64).$(EXTENSION_BIN_LINUX)
BIN_LINUX_ARM64 := $(BIN_DIR)/$(service).$(GOARCH_ARM64).$(EXTENSION_BIN_LINUX)
BIN_DARWIN_AMD64 := $(BIN_DIR)/$(service).$(GOARCH_AMD64).$(EXTENSION_BIN_DARWIN)
BIN_DARWIN_ARM64 := $(BIN_DIR)/$(service).$(GOARCH_ARM64).$(EXTENSION_BIN_DARWIN)
BIN_WINDOWS_AMD64 := $(BIN_DIR)/$(service).$(GOARCH_AMD64).$(EXTENSION_BIN_WINDOWS)
BIN_WINDOWS_ARM64 := $(BIN_DIR)/$(service).$(GOARCH_ARM64).$(EXTENSION_BIN_WINDOWS)

define go_build_platform
	env GOOS=$(1) GOARCH=$(2) go build -o $(3) $(SELECTED_SERVICE_DIR)
endef

.PHONY: list-microservices
list-microservices:
	ls $(CMD_DIR)

.PHONY: wire
wire:
	cd $(SELECTED_SERVICE_DIR); \
	go run github.com/google/wire/cmd/wire gen $(SELECTED_SERVICE_ROOT_PKG)/app

.PHONY: run
run: wire
	go run $(SELECTED_SERVICE_DIR)

.PHONY: build-linux-amd64
build-linux-amd64: wire
	$(call go_build_platform,$(GOOS_LINUX),$(GOARCH_AMD64),$(BIN_LINUX_AMD64))

.PHONY: build-linux-arm64
build-linux-arm64: wire
	$(call go_build_platform,$(GOOS_LINUX),$(GOARCH_ARM64),$(BIN_LINUX_ARM64))

.PHONY: build-darwin-amd64
build-darwin-amd64: wire
	$(call go_build_platform,$(GOOS_DARWIN),$(GOARCH_AMD64),$(BIN_DARWIN_AMD64))

.PHONY: build-darwin-arm64
build-darwin-arm64: wire
	$(call go_build_platform,$(GOOS_DARWIN),$(GOARCH_ARM64),$(BIN_DARWIN_ARM64))

.PHONY: build-windows-amd64
build-windows-amd64: wire
	$(call go_build_platform,$(GOOS_WINDOWS),$(GOARCH_AMD64),$(BIN_WINDOWS_AMD64))

.PHONY: build-windows-arm64
build-windows-arm64: wire
	$(call go_build_platform,$(GOOS_WINDOWS),$(GOARCH_ARM64),$(BIN_WINDOWS_ARM64))

.PHONY: build-all-platforms
build-all-platforms: wire build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-windows-arm64

.PHONY: run-bin-linux-amd64
run-bin-linux-amd64: build-linux-amd64
	exec $(BIN_LINUX_AMD64)

.PHONY: run-bin-linux-arm64
run-bin-linux-arm64: build-linux-arm64
	exec $(BIN_LINUX_ARM64)

.PHONY: run-bin-darwin-amd64
run-bin-darwin-amd64: build-darwin-amd64
	exec $(BIN_DARWIN_AMD64)

.PHONY: run-bin-darwin-arm64
run-bin-darwin-arm64: build-darwin-arm64
	exec $(BIN_DARWIN_ARM64)

.PHONY: run-bin-windows-amd64
run-bin-windows-amd64: build-windows-amd64
	exec $(BIN_WINDOWS_AMD64)

.PHONY: run-bin-windows-arm64
run-bin-windows-arm64: build-windows-arm64
	exec $(BIN_WINDOWS_ARM64)

.PHONY: test
test: wire
ifeq ($(testlocation),)
	go test $(ROOT_DIR)/...
else ifeq ($(suffix $(TEST_PATH)), .go)
	go test $(TEST_PATH)
else
	go test $(TEST_PATH)/...
endif

clean:
	rm -rf $(BIN_DIR)/*