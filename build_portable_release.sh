#!/bin/bash

# Check for version argument
if [ -z "$1" ]; then
  echo "Usage: $0 <version> (e.g., $0 v1.0.0)"
  exit 1
fi

# Configuration
VERSION="$1"
APP_NAME="GameHub"
BINARY_BASE_NAME="game_hub"
MODULE_PATH="game_hub/core"
PLATFORMS=("linux/amd64" "linux/arm64" "windows/amd64" "windows/arm64" "darwin/amd64" "darwin/arm64")
RELEASE_DIR="release/portable"

# Check for 7z availability
USE_7Z=true
if ! command -v 7z &> /dev/null; then
  echo "Note: 7z not found, will use zip (slower)"
  USE_7Z=false
fi

# Clean and prepare release directory
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS=${PLATFORM%%/*}
  GOARCH=${PLATFORM##*/}

  # Set binary name
  BINARY_NAME="$BINARY_BASE_NAME"
  [ "$GOOS" = "windows" ] && BINARY_NAME="${BINARY_NAME}.exe"

  # Prepare output directories
  OUTPUT_DIR="${RELEASE_DIR}/${BINARY_BASE_NAME}_${GOOS}_${GOARCH}_${VERSION}"
  DATA_DIR="${OUTPUT_DIR}/data"
  mkdir -p "$OUTPUT_DIR" "$DATA_DIR"

  echo "Building for $GOOS/$GOARCH..."
  cd game_hub
  BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X ${MODULE_PATH}.AppName=${APP_NAME} -X ${MODULE_PATH}.Version=${VERSION} -X ${MODULE_PATH}.BuildTime=${BUILD_TIME}" -o "../${OUTPUT_DIR}/${BINARY_NAME}"
  cd ../

  # Copy data files if they exist
  if [ -d "./data" ]; then
    cp -r ./data/* "$DATA_DIR/"
  fi

  cd "$RELEASE_DIR"
  ARCHIVE_NAME="${BINARY_BASE_NAME}_${GOOS}_${GOARCH}_${VERSION}"
  
  if [ "$GOOS" = "windows" ]; then
    if [ "$USE_7Z" = true ]; then
      echo "Creating 7z archive..."
      7z a -mx=9 "${ARCHIVE_NAME}.7z" "$ARCHIVE_NAME" > /dev/null 2>&1
    else
      echo "Creating zip archive..."
      zip -qr "${ARCHIVE_NAME}.zip" "$ARCHIVE_NAME" > /dev/null 2>&1
    fi
  else
    echo "Creating tar.gz archive..."
    tar -czf "${ARCHIVE_NAME}.tar.gz" "$ARCHIVE_NAME"
  fi
  
  #rm -rf "$ARCHIVE_NAME"
  cd - > /dev/null
done

echo "Release built in ${RELEASE_DIR}/ directory."