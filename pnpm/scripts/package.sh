#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Define variables for directories
REPO_DIR=$(dirname $(dirname "$0"))
BUILD_DIR="$REPO_DIR/build"
BIN_DIR="$REPO_DIR/bin"
TARBALL_NAME="pnpm-buildpack.tar.gz"
CNB_NAME="pnpm-buildpack.cnb"

# Create the build directory if it doesn't exist
mkdir -p "$BUILD_DIR"

# Package the buildpack files into a tarball
echo "Packaging buildpack into tarball..."
tar -czf "$BUILD_DIR/$TARBALL_NAME" -C "$REPO_DIR" bin buildpack.toml

# Package the tarball into a .cnb file
echo "Packaging tarball into .cnb file..."
mv "$BUILD_DIR/$TARBALL_NAME" "$BUILD_DIR/$CNB_NAME"

echo "Packaging completed. Artifacts are in the build directory."
