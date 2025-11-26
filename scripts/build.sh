#!/usr/bin/env bash
set -euo pipefail

BINARY="flint"
SRC="./cmd/flint"
DIST="./dist"
VERSION="${VERSION:-0.1.0.0}-"
BUILD_TYPE="alpha"
FULL_VERSION="${VERSION}${BUILD_TYPE}"
GIT_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
LDFLAGS="-X 'flint/internal/version.Version=$FULL_VERSION'
-X 'flint/internal/version.GitCommit=$GIT_COMMIT'
-X 'flint/internal/version.BuildTime=$BUILD_TIME'"

TARGETS=(
    "linux amd64"
    "linux arm64"
    "linux arm"
    "darwin amd64"
    "darwin arm64"
    "windows amd64"
    "windows arm64"
)

mkdir -p "$DIST"

for TARGET in "${TARGETS[@]}"; do
    read -r GOOS GOARCH <<< "$TARGET"
    OUT="$DIST/${BINARY}-${GOOS}-${GOARCH}"
    [[ "$GOOS" == "windows" ]] && OUT+=".exe"
    echo "Building $OUT ..."
    GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags "$LDFLAGS" -o "$OUT" "$SRC"
done

echo "Build finished! Binaries are in $DIST/"