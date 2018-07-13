## Folder content generated files
BUILD_FOLDER = ./build

## command
GO           = go
GO_VENDOR    = govendor
MKDIR_P      = mkdir -p
BATS         = bats
DOCKER       = docker

################################################

.PHONY: all
all: build test

.PHONY: pre-build
pre-build: govendor-sync

.PHONY: build
build: pre-build src.build

.PHONY: test
test: build src.test src.cmd.bats

.PHONY: check
check: check-govendor check-bats check-docker

.PHONY: clean
clean:
	$(RM) -rf $(BUILD_FOLDER)
	$(GO) clean -i -r -x -cache -testcache

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
	@[ "`which $(GO_VENDOR)`" != "" ] || (echo "$(GO_VENDOR) is missing"; false) && (echo ".. OK")

.PHONY: check-bats
check-bats:
	$(info check bats)
	@[ "`which $(BATS)`" != "" ] || (echo "$(BATS) is missing"; false) && (echo ".. OK")

.PHONY: check-docker
check-docker:
	$(info check docker)
	@[ "`which $(DOCKER)`" != "" ] || (echo "$(DOCKER) is missing"; false) && (echo ".. OK")

## docker #######################################

.PHONY: docker.build
docker.build:
	$(DOCKER) build --tag linkernetworks/oauth:latest .
