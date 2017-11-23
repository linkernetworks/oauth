GO_BUILD_FLAGS=
GO_TEST_FLAGS=

SHELL := /bin/bash

PKGS=$(shell go list ./... | grep -E -v "(vendor|cmd)")

all:
	go build $(GO_BUILD_FLAG) ./cmd/lnk-auth

install:
	go install $(GO_BUILD_FLAG) ./cmd/lnk-auth

clean:
	rm -f lnk-auth

test:
	# refer to https://github.com/golang/go/issues/11659#issuecomment-122139338
	go test $(GO_TEST_FLAGS) --cover $(PKGS)

test-debug:
	# refer to https://github.com/golang/go/issues/11659#issuecomment-122139338
	go test $(GO_TEST_FLAGS) --cover $(PKGS) -test.v -x

dev:
	./lnk-auth --config config/dev.json
