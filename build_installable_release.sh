#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <version> (e.g., $0 v1.0.0)"
  exit 1
fi

VERSION="$1"
APP_NAME="GameHub"
BINARY_BASE_NAME="game_hub"
MODULE_PATH="game_hub/core"
INSTALLER_MODULE_PATH="main"
PLATFORMS=("linux/amd64" "linux/386" "linux/arm64" "windows/amd64" "windows/arm64" "darwin/amd64" "darwin/arm64")
RELEASE_DIR="release/installers"

USE_7Z=true
if ! command -v 7z &> /dev/null; then
  echo "Note: 7z not found, will use zip (slower)"
  USE_7Z=false
fi

rm -rf installer/data
mkdir -p installer/data
if [ -d "./data" ]; then
  cp -r ./data/* installer/data/
fi

mkdir -p installer/install

rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS=${PLATFORM%%/*}
  GOARCH=${PLATFORM##*/}

  GAME_BINARY_NAME="$BINARY_BASE_NAME"
  INSTALLER_BINARY_NAME="${BINARY_BASE_NAME}_installer"
  if [ "$GOOS" = "windows" ]; then
    GAME_BINARY_NAME="${GAME_BINARY_NAME}.exe"
    INSTALLER_BINARY_NAME="${INSTALLER_BINARY_NAME}.exe"
  fi

  rm -rf installer/install/*
  mkdir -p installer/install

  cd game_hub
  echo "Building app binary for $GOOS/$GOARCH..."
  BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X ${MODULE_PATH}.AppName=${APP_NAME} -X ${MODULE_PATH}.Version=${VERSION} -X ${MODULE_PATH}.BuildTime=${BUILD_TIME}" -o "../installer/install/${GAME_BINARY_NAME}"
  cd ../

  if [ "$GOOS" = "linux" ] || [ "$GOOS" = "darwin" ]; then
    cd installer
    echo "Generating game wrapper for $GOOS..."
    ./generate_wrapper.sh "$GOOS" "$APP_NAME" "$VERSION" "../installer/install/$GAME_BINARY_NAME" "../installer/install"
    cd ../
  fi
  if [ "$GOOS" = "darwin" ]; then
    rm installer/install/$GAME_BINARY_NAME
  fi

  OUTPUT_DIR="${RELEASE_DIR}/${BINARY_BASE_NAME}_installer_${GOOS}_${GOARCH}_${VERSION}"
  mkdir -p "$OUTPUT_DIR"

  cd installer
  echo "Building installer for $GOOS/$GOARCH..."
  INSTALLER_OUTPUT="../${OUTPUT_DIR}/${INSTALLER_BINARY_NAME}"
  GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X ${INSTALLER_MODULE_PATH}.AppName=${APP_NAME} -X ${INSTALLER_MODULE_PATH}.BinaryName=${GAME_BINARY_NAME}" -o "$INSTALLER_OUTPUT"
  cd ../

  if [ "$GOOS" = "darwin" ]; then
    cd installer
    echo "Generating installer wrapper for $GOOS..."
    ./generate_wrapper.sh "$GOOS" "${APP_NAME}Installer" "$VERSION" "$INSTALLER_OUTPUT" "../$OUTPUT_DIR"
    rm "$INSTALLER_OUTPUT"
    cd ../
  fi

  cd "$RELEASE_DIR"
  ARCHIVE_NAME="${BINARY_BASE_NAME}_installer_${GOOS}_${GOARCH}_${VERSION}"
  
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

echo "Installers built in ${RELEASE_DIR}/ directory."