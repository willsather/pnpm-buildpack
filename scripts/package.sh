#!/bin/bash

# directory locations
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# import utils
source "${SCRIPT_DIR}/utils/print.sh"

# define variables
BUILD_DIR="bin"
PACKAGE_NAME="buildpack.tar.gz"
BUILDPACK_TOML="buildpack.toml"

util::print::title "** Buildpack packaging **"

# check if build directory exists
if [ ! -d "$BUILD_DIR" ]; then
    echo "Error: Build directory '$BUILD_DIR' does not exist."
    exit 1
fi

# check if buildpack.toml exists
if [ ! -f "$BUILDPACK_TOML" ]; then
    echo "Error: '$BUILDPACK_TOML' file does not exist."
    exit 1
fi

# create the tarball
echo "- creating the package $PACKAGE_NAME..."
tar -czvf $PACKAGE_NAME $BUILDPACK_TOML $BUILD_DIR

# check if packaging status
if [ $? -eq 0 ]; then
    util::print::success "** Buildpack packaging completed. Saved as $PACKAGE_NAME **"
else
    util::print::error "** Buildpack packaging failed **"
    exit 1
fi
