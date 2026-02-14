#!/usr/bin/env bash
set -euo pipefail

echo "[verify] running go tests..."
go test ./...

echo "[verify] running API smoke checks..."
/app/scripts/smoke_api.sh

echo "[verify] all checks passed"
