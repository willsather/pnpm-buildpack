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
    util::print::info "... Building buildpack './pnpm/build/pnpm.cnb'"
    "$PNPM_BUILD_SCRIPT"
  else
    util::print::error "Error: $PNPM_BUILD_SCRIPT not found or not executable."
  fi

  if [ -x "$PNPM_INSTALL_BUILD_SCRIPT" ]; then
    util::print::info "... Building buildpack './pnpm-install/build/pnpm-install.cnb'"
    "$PNPM_INSTALL_BUILD_SCRIPT"
  else
    util::print::error "Error: $PNPM_INSTALL_BUILD_SCRIPT not found or not executable."
  fi

  if [ -x "$PNPM_START_BUILD_SCRIPT" ]; then
    util::print::info "... Building buildpack './pnpm-start/build/pnpm-start.cnb'"
    "$PNPM_START_BUILD_SCRIPT"
  else
    util::print::error "Error: $PNPM_START_BUILD_SCRIPT not found or not executable."
  fi

  if [ $? -eq 0 ]; then
    util::print::success "** Buildpack packaging completed **"
  else
    util::print::error "** Buildpack packaging failed **"
  fi
}

main "$@"
