#!/usr/bin/env bash

set -eu
set -o pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILD_SCRIPT="${SCRIPT_DIR}/build.sh"
readonly PACKAGE_SCRIPT="${SCRIPT_DIR}/package.sh"

source "${SCRIPT_DIR}/utils/print.sh"

function main() {
  # Check if build.sh exists and is executable
  if [[ -x "${BUILD_SCRIPT}" ]]; then
    "${BUILD_SCRIPT}"
  else
    util::print::error "Error: Build script '${BUILD_SCRIPT}' does not exist or is not executable."
  fi

  # Check if package.sh exists and is executable
  if [[ -x "${PACKAGE_SCRIPT}" ]]; then
    "${PACKAGE_SCRIPT}"
  else
    util::print::error "Error: Package script '${PACKAGE_SCRIPT}' does not exist or is not executable."
  fi
}

main "$@"
