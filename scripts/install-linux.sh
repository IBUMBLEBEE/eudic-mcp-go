#!/usr/bin/env bash
# Install eudic-mcp-go into ~/bin on a Linux host (Cursor Remote SSH).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
mkdir -p "${HOME}/bin"

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  SRC="$ROOT/dist/linux-amd64/eudic-mcp-go" ;;
  aarch64|arm64) SRC="$ROOT/dist/linux-arm64/eudic-mcp-go" ;;
  *)
    echo "Unsupported arch: $ARCH — build locally with: go build -o eudic-mcp-go ."
    exit 1
    ;;
esac

if [[ ! -f "$SRC" ]]; then
  echo "Missing $SRC — run scripts/build.sh first, or:"
  echo "  go build -ldflags='-s -w' -o ${HOME}/bin/eudic-mcp-go $ROOT"
  if command -v go >/dev/null 2>&1; then
    go build -ldflags="-s -w" -o "${HOME}/bin/eudic-mcp-go" "$ROOT"
    chmod +x "${HOME}/bin/eudic-mcp-go"
    echo "Built and installed to ${HOME}/bin/eudic-mcp-go"
    exit 0
  fi
  exit 1
fi

cp "$SRC" "${HOME}/bin/eudic-mcp-go"
chmod +x "${HOME}/bin/eudic-mcp-go"
echo "Installed ${HOME}/bin/eudic-mcp-go"
echo "Configure remote ~/.cursor/mcp.json command to that path, then Reload Window."
