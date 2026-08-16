#!/usr/bin/env bash
# Build eudic-mcp-go for common platforms (run on any machine with Go).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p dist/linux-amd64 dist/linux-arm64 dist/windows-amd64
export CGO_ENABLED=0

echo "→ linux/amd64"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/linux-amd64/eudic-mcp-go .

echo "→ linux/arm64"
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/linux-arm64/eudic-mcp-go .

echo "→ windows/amd64"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/windows-amd64/eudic-mcp-go.exe .

# Local host binary (for non-cross day-to-day use)
go build -ldflags="-s -w" -o eudic-mcp-go .

echo "Done. Artifacts under dist/"
ls -la dist/*/eudic-mcp-go* 2>/dev/null || true
