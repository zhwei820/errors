#!/usr/bin/env bash
set -e

# Install protoc
#
#   - Mac OS: `brew install protobuf`
#   - Ubuntu: sudo apt install protobuf-compiler
#
#  go install google.golang.org/protobuf/cmd/protoc-gen-go
#  go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Get the directory of this script file.
DIR=$(dirname "$0")

protoc -I=. --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative ./*.proto
