#!/bin/bash

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "${SCRIPT_DIR}/utils/print.sh"

BUILD_OUTPUT_DIR="./build"
PACKAGE_CONFIG="./package.toml"

function main() {
  util::print::title "** Buildpack packaging **"

  mkdir -p "$BUILD_OUTPUT_DIR"

  util::print::info "... Packaging buildpack into $BUILD_OUTPUT_DIR/pnpm-buildpack.cnb ..."
  pack buildpack package "$BUILD_OUTPUT_DIR/pnpm-buildpack.cnb" --config "$PACKAGE_CONFIG" --format file

  if [ $? -eq 0 ]; then
    util::print::success "** Buildpack packaging completed **"
  else
    util::print::error "** Buildpack packaging failed **"
  fi
}

main "$@"
