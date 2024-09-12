#!/usr/bin/env bash

set -eu
set -o pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACK_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

source "${SCRIPT_DIR}/utils/print.sh"

function main() {
  util::print::title "Run Buildpack Unit Tests"

    pushd "${BUILDPACK_DIR}" > /dev/null
      if go test ./... -v -short; then # -short excludes integration tests
        util::print::success "** GO Test Succeeded **"
      else
        util::print::error "** GO Test Failed **"
      fi
    popd > /dev/null
}

main "${@:-}"
