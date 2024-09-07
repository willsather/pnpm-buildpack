#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Define variables for directories
REPO_DIR=$(dirname $(dirname "$0"))
BUILD_DIR="$REPO_DIR/build"
BIN_DIR="$REPO_DIR/bin"
TARBALL_NAME="pnpm-start-buildpack.tar.gz"
CNB_NAME="pnpm-start-buildpack.cnb"

# Create the build directory if it doesn't exist
mkdir -p "$BUILD_DIR"

# Package the buildpack files into a tarball
echo "Packaging buildpack into tarball..."
tar -czf "$BUILD_DIR/$TARBALL_NAME" -C "$REPO_DIR" bin buildpack.toml

# Create the .cnb file by copying the tarball
echo "Creating .cnb file from the tarball..."
cp "$BUILD_DIR/$TARBALL_NAME" "$BUILD_DIR/$CNB_NAME"

echo "Packaging completed. Artifacts (tarball and .cnb file) are in the build directory."
