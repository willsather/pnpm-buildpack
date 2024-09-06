#!/usr/bin/env bash

set -eu
set -o pipefail

# Directory locations
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILD_SCRIPT="${SCRIPT_DIR}/go_build.sh"
readonly PACKAGE_SCRIPT="${SCRIPT_DIR}/package.sh"

function main() {
  # Check if build.sh exists and is executable
  if [[ -x "${BUILD_SCRIPT}" ]]; then
    "${BUILD_SCRIPT}"
  else
    util::print::error "Error: Build script '${BUILD_SCRIPT}' does not exist or is not executable."
    exit 1
  fi

  # Check if package.sh exists and is executable
  if [[ -x "${PACKAGE_SCRIPT}" ]]; then
    "${PACKAGE_SCRIPT}"
  else
    util::print::error "Error: Package script '${PACKAGE_SCRIPT}' does not exist or is not executable."
    exit 1
  fi
}

main "$@"
