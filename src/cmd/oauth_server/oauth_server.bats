#!/usr/bin/env bats

#
# How to run this test:
#   Install bats: https://github.com/bats-core/bats-core#installation
#   Run: bats [file name]
#

setup() {
  ## create a temporary folder for this test

  ## build oauth_server binary
  binaryFile="${BATS_TMPDIR}/oauth_server"
  go build -o ${binaryFile} ./src/cmd/oauth_server/main.go/..
}

## teardown(){}

@test "oauth_server can be killed by TERM signal" {
  ## arrange
  $binaryFile &
  pid=$!
  echo pid = $pid

  ## action: sleep a while and then kill pid in background
  (sleep 1 && kill -s TERM $pid) &

  ## assert
  wait $pid
}

@test "oauth_server can be killed by INT signal" {
  ## arrange
  $binaryFile &
  pid=$!
  echo pid = $pid

  ## action: sleep a while and then kill pid in background
  (sleep 1 && kill -s INT $pid) &

  ## assert
  wait $pid
}
