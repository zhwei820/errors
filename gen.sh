#!/usr/bin/env bash
set -e

# Script that generates .pb.go files from the .proto files.
# Manually re-run this script after you edit any .proto file.

# 1. Install protoc
#
#   - Mac OS: `brew install protobuf`
#   - Ubuntu: sudo apt install protobuf-compiler
#
# 2. Install or update gRPC and the protoc plugin for Golang if you have not:
#
# $ go install \
#     google.golang.org/protobuf/cmd/protoc-gen-go \
#     google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Get the directory of this script file.
DIR=$(dirname "$0")

protoc -I=. \
   --go_out . --go_opt paths=source_relative \
   --go-grpc_out . --go-grpc_opt paths=source_relative \
   ./*.proto

