#!/bin/bash

# directory locations
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACK_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
readonly BIN_DIR="bin"
readonly BUILD_DIR="${BUILDPACK_DIR}/build"

# import utils
source "${SCRIPT_DIR}/utils/print.sh"

# define variables
PACKAGE_NAME="pnpm-install-buildpack.tar.gz"
BUILDPACK_TOML="buildpack.toml"
CNB_PACKAGE_NAME="pnpm-install-buildpack.cnb"

function main() {
    util::print::title "** Buildpack packaging **"

    if [ ! -d "${BUILD_DIR}" ]; then
        echo "- Creating the output directory ${BUILD_DIR}..."
        mkdir -p "${BUILD_DIR}"
    fi

    package

    util::print::success "** Buildpack packaging completed **"
}

function package() {
    util::print::title "** Packaging Buildpack **"

    # check if buildpack.toml exists
    if [ ! -f "$BUILDPACK_TOML" ]; then
        util::print::error "Error: '$BUILDPACK_TOML' file does not exist."
        exit 1
    fi

    # create the tarball
    echo "- Creating the package $PACKAGE_NAME..."
    tar -czvf "${BUILD_DIR}/${PACKAGE_NAME}" "$BUILDPACK_TOML" "$BIN_DIR"

    # check if tarball packaging status
    if [ $? -ne 0 ]; then
        util::print::error "** Buildpack TAR.GZ packaging failed **"
        exit 1
    fi

    # create the .cnb file
    echo "- Creating the .cnb package $CNB_PACKAGE_NAME..."
    cp "${BUILD_DIR}/${PACKAGE_NAME}" "${BUILD_DIR}/${CNB_PACKAGE_NAME}"

    # check if .cnb packaging status
    if [ $? -ne 0 ]; then
        util::print::error "** Buildpack CNB packaging failed **"
        exit 1
    fi
}

# Execute the functions
main "$@"
