#!/usr/bin/env bash
# Build one MCPB bundle per GOOS/GOARCH target.
# Usage: scripts/pack-mcpb.sh [version]
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

VERSION="${1:-${MCPB_VERSION:-}}"
if [[ -z "$VERSION" ]]; then
  if git describe --tags --exact-match >/dev/null 2>&1; then
    VERSION="$(git describe --tags --exact-match)"
  elif git describe --tags --abbrev=0 >/dev/null 2>&1; then
    VERSION="$(git describe --tags --abbrev=0)"
  else
    VERSION="$(node -e 'console.log(require("./mcpb/manifest.json").version)')"
  fi
fi
VERSION="${VERSION#v}"

OUT_DIR="${MCPB_OUT_DIR:-$ROOT/mcpb-dist}"
STAGE_ROOT="$OUT_DIR/staging"
mkdir -p "$OUT_DIR"
rm -rf "$STAGE_ROOT"

mcpb_cmd() {
  if command -v mcpb >/dev/null 2>&1; then
    mcpb "$@"
  else
    npx --yes @anthropic-ai/mcpb "$@"
  fi
}

mcpb_cmd validate "$ROOT/mcpb/manifest.json"

targets=(
  "darwin arm64 darwin darwin-arm64"
  "darwin amd64 darwin darwin-amd64"
  "linux amd64 linux linux-amd64"
  "linux arm64 linux linux-arm64"
  "windows amd64 win32 windows-amd64"
  "windows arm64 win32 windows-arm64"
)

for spec in "${targets[@]}"; do
  read -r goos goarch platform suffix <<<"$spec"

  bin_name="quark-mcp"
  if [[ "$goos" == "windows" ]]; then
    bin_name="quark-mcp.exe"
  fi

  stage="$STAGE_ROOT/$suffix"
  mkdir -p "$stage/server"

  echo "Building $goos/$goarch"
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" go build -trimpath \
    -ldflags="-s -w -X github.com/chiehw/quark-mcp/internal/version.Version=${VERSION}" \
    -o "$stage/server/$bin_name" .
  chmod 0755 "$stage/server/$bin_name"

  node "$ROOT/scripts/stamp-mcpb-manifest.mjs" \
    "$ROOT/mcpb/manifest.json" \
    "$stage/manifest.json" \
    "$VERSION" \
    "$platform" \
    "$bin_name"

  mcpb_cmd validate "$stage/manifest.json"

  output="$OUT_DIR/quark-mcp_${VERSION}_${suffix}.mcpb"
  rm -f "$output"
  mcpb_cmd pack "$stage" "$output"
  mcpb_cmd info "$output"
done

rm -rf "$STAGE_ROOT"
shopt -s nullglob
bundles=("$OUT_DIR"/quark-mcp_"${VERSION}"_*.mcpb)
if [[ ${#bundles[@]} -ne 6 ]]; then
  echo "expected 6 MCPB bundles, found ${#bundles[@]}" >&2
  ls -1 "$OUT_DIR"/*.mcpb >&2 || true
  exit 1
fi
echo "MCPB bundles written to $OUT_DIR"
ls -1 "$OUT_DIR"/*.mcpb
