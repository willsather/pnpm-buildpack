#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Define variables for directories
REPO_DIR=$(dirname $(dirname "$0"))
BIN_DIR="$REPO_DIR/bin"

# Create the bin directory if it doesn't exist
mkdir -p "$BIN_DIR"

# Build the detect and build binaries
echo "Building detect binary..."
go build -o "$BIN_DIR/detect" "$REPO_DIR/run/main.go"

echo "Building build binary..."
go build -o "$BIN_DIR/build" "$REPO_DIR/run/main.go"

echo "Build completed. Binaries are in the bin directory."
