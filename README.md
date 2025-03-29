# gRPC Learning Repository

This repository is for practicing the gRPC learning content introduced in the following article:

[Introduction to gRPC: API Development with Protocol Buffers and gRPC](https://note.com/shunex/n/nd8109a1144a5)

## Overview

gRPC is a high-performance RPC framework developed by Google that uses Protocol Buffers (protobuf) for efficient communication. In this repository, you will learn and practice the following:

- Defining data using Protocol Buffers
- Implementing gRPC servers and clients
- Using gRPC for bidirectional streaming communication

## Setup

1. Install the required tools:
   - `protoc` (Protocol Buffers compiler)
   - `protoc-gen-go` and `protoc-gen-go-grpc` (Code generation plugins for Go)

   ```bash
   brew install protobuf
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

2. Install project dependencies:

   ```bash
   go mod tidy
   ```

3. Generate code from Protocol Buffers files:

   ```bash
   protoc --go_out=. --go-grpc_out=. chat/chat.proto
   ```

## How to Run

1. Start the gRPC server:

   ```bash
   go run server/main.go
   ```

2. Run the gRPC client:

   ```bash
   go run client/main.go
   ```

## References

- [gRPC Official Documentation](https://grpc.io/docs/)
- [Protocol Buffers Official Documentation](https://protobuf.dev/)

---

Let's dive into gRPC and master it from basics to advanced! 🚀✨
