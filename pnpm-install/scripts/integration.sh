#!/usr/bin/env bash

set -eu
set -o pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACK_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

source "${SCRIPT_DIR}/utils/print.sh"

function main() {
  if [[ ! -d "${BUILDPACK_DIR}/integration" ]]; then
      util::print::warn "** WARNING  No Integration tests **"
  fi


#  TODO: Should this script install and cycle through builders?
#  util::print::info "- Setting default pack builder to 'paketobuildpacks/builder-jammy-base'"
#  pack config default-builder paketobuildpacks/builder-jammy-base

  util::print::title "Run Buildpack Runtime Integration Tests"

  export CGO_ENABLED=0
  pushd "${BUILDPACK_DIR}" > /dev/null
    if go test ./... -v; then
      util::print::info "** GO Test Succeeded **"
    else
      util::print::error "** GO Test Failed **"
    fi
  popd > /dev/null

  util::print::success "** GO Test Succeeded**"
}


main "${@:-}"