#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
  --config "$ROOT_DIR/files/codegen/server.cfg.yaml" \
  --package http \
  -o "$ROOT_DIR/internal/controller/http/openapi_server.gen.go" \
  "$ROOT_DIR/internal/controller/http/contract.yaml"
