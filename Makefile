## Folder content generated files
BUILD_FOLDER = ./build

## command
GO           = go
GO_VENDOR    = govendor
MKDIR_P      = mkdir -p
BATS         = bats

################################################

.PHONY: all
all: build test

.PHONY: pre-build
pre-build:
	$(MAKE) govendor-sync

.PHONY: build
build: pre-build
	$(MAKE) src.build

.PHONY: test
test: build
	$(MAKE) src.test
	$(MAKE) src.cmd.bats

.PHONY: check
check:
	$(MAKE) check-govendor

.PHONY: clean
clean:
	$(RM) -rf $(BUILD_FOLDER)

## vendor/ #####################################

.PHONY: govendor-sync
govendor-sync:
	$(GO_VENDOR) sync -v
	$(GO_VENDOR) remove -v +unused

## src/ ########################################

.PHONY: src.build
src.build:
	$(GO) build -v ./src/...
	$(MKDIR_P) $(BUILD_FOLDER)/src/cmd/lnk-auth/
	$(GO) build -v -o $(BUILD_FOLDER)/src/cmd/oauth_server/oauth_server ./src/cmd/oauth_server/...

.PHONY: src.test
src.test:
	$(GO) test -v -race ./src/...

.PHONY: src.install
src.install:
	$(GO) install -v ./src/...

.PHONY: src.test-coverage
src.test-coverage:
	$(MKDIR_P) $(BUILD_FOLDER)/src/
	$(GO) test -v -race -coverprofile=$(BUILD_FOLDER)/src/coverage.txt -covermode=atomic ./src/...
	$(GO) tool cover -html=$(BUILD_FOLDER)/src/coverage.txt -o $(BUILD_FOLDER)/src/coverage.html

## src/cmd/ ####################################

.PHONY: src.cmd.bats
src.cmd.bats:
	@$(BATS) -t $(shell find src/cmd -name "*.bats")

## check build env #############################

.PHONY: check-govendor
check-govendor:
	$(info check govendor)
	@[ "`which $(GO_VENDOR)`" != "" ] || (echo "$(GO_VENDOR) is missing"; false)

.PHONY: check-bats
check-bats:
	$(info check bats)
	@[ "`which $(BATS)`" != "" ] || (echo "$(BATS) is missing"; false)

## docker #######################################

.PHONY: dockerfiles.build
docker.build:
	docker build --tag linkernetworks/oauth:latest .
