#!/bin/bash

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

source "${SCRIPT_DIR}/utils/print.sh"

PNPM_BUILD_SCRIPT="$ROOT_DIR/pnpm/scripts/build.sh"
PNPM_INSTALL_BUILD_SCRIPT="$ROOT_DIR/pnpm-install/scripts/build.sh"
PNPM_START_BUILD_SCRIPT="$ROOT_DIR/pnpm-start/scripts/build.sh"

function main() {
  util::print::title "** Building buildpacks **"

  if [ -x "$PNPM_BUILD_SCRIPT" ]; then
    "$PNPM_BUILD_SCRIPT"
  else
    util::print::error "Error: $PNPM_BUILD_SCRIPT not found or not executable."
  fi

  if [ -x "$PNPM_INSTALL_BUILD_SCRIPT" ]; then
    "$PNPM_INSTALL_BUILD_SCRIPT"
  else
    util::print::error "Error: $PNPM_INSTALL_BUILD_SCRIPT not found or not executable."
  fi

  if [ -x "$PNPM_START_BUILD_SCRIPT" ]; then
    "$PNPM_START_BUILD_SCRIPT"
  else
    util::print::error "Error: $PNPM_START_BUILD_SCRIPT not found or not executable."
  fi

  if [ $? -eq 0 ]; then
    echo
    echo
    util::print::success "** Buildpacks build completed **"
  else
    util::print::error "** Buildpack packaging failed **"
  fi
}

main "$@"
