//go:build tools
// +build tools

package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/sqlc-dev/sqlc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

// This file pins versions of build-time tools (e.g. protoc, sqlc)
// in go.mod to ensure reproducible code generation across environments (local, CI, production).
