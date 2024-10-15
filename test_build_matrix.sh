#!/bin/bash

set -e

GO_VERSIONS=("1.20")
OS_LIST=("windows" "darwin" "linux")
ARCH_LIST=("amd64" "arm64")

# Run tests only once on the current platform
echo "======================================"
echo "Running tests on your current platform"
echo "======================================"

go test -v ./...

for GO_VERSION in "${GO_VERSIONS[@]}"; do
  for OS in "${OS_LIST[@]}"; do
    for ARCH in "${ARCH_LIST[@]}"; do
      # Skip unsupported combinations
      if [[ "$OS" == "windows" && "$ARCH" == "arm64" ]]; then
        continue
      fi

      echo "======================================"
      echo "Building for Go $GO_VERSION on $OS $ARCH"
      echo "======================================"

      export GOOS=$OS
      export GOARCH=$ARCH

      # Build binary
      OUTPUT_NAME="cleaner-$GOOS-$GOARCH"
      if [ "$GOOS" == "windows" ]; then
        OUTPUT_NAME+=".exe"
      fi

      echo "Building binary..."
      go build -ldflags="-s -w" -o build/$OUTPUT_NAME

      echo "Built $OUTPUT_NAME"
      echo
    done
  done
done
