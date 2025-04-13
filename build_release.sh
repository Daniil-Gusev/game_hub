#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <version> (e.g., $0 v1.0.0)"
  exit 1
fi

VERSION="$1"

PLATFORMS=("linux/amd64" "linux/arm64" "windows/amd64" "windows/arm64" "darwin/amd64" "darwin/arm64")

rm -rf release
mkdir -p release

for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS=${PLATFORM%%/*}
  GOARCH=${PLATFORM##*/}

  BINARY_NAME="game_hub"
  if [ "$GOOS" = "windows" ]; then
    BINARY_NAME="game_hub.exe"
  fi

  OUTPUT_DIR="release/game_hub_${GOOS}_${GOARCH}_${VERSION}"
  mkdir -p "$OUTPUT_DIR"

  echo "Building for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/$BINARY_NAME" .

  mkdir -p "$OUTPUT_DIR/core" "$OUTPUT_DIR/app" "$OUTPUT_DIR/games"
  cp core/*.json "$OUTPUT_DIR/core/" 2>/dev/null || true
  cp app/*.json "$OUTPUT_DIR/app/" 2>/dev/null || true
  cp games/*.json "$OUTPUT_DIR/games/" 2>/dev/null || true

  find games -type d -not -path "games/game_template*" | while read -r dir; do
    if [ "$dir" = "games" ] || [ "$dir" = "games/game_template" ]; then
      continue
    fi
    rel_dir=${dir#games/}
    mkdir -p "$OUTPUT_DIR/games/$rel_dir"
    cp "$dir"/*.json "$OUTPUT_DIR/games/$rel_dir/" 2>/dev/null || true
  done

  cd release
  if [ "$GOOS" = "windows" ]; then
    zip -r "game_hub_${GOOS}_${GOARCH}_${VERSION}.zip" "game_hub_${GOOS}_${GOARCH}_${VERSION}"
  else
    tar -czf "game_hub_${GOOS}_${GOARCH}_${VERSION}.tar.gz" "game_hub_${GOOS}_${GOARCH}_${VERSION}"
  fi
  rm -rf "game_hub_${GOOS}_${GOARCH}_${VERSION}"
  cd ..
done

echo "Release built in release/ directory."