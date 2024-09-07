#!/usr/bin/env bash

set -eu
set -o pipefail

# Directory locations
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACK_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

# import utils
source "${SCRIPT_DIR}/utils/print.sh"

function main() {
  util::print::title "** GO build **"

  mkdir -p "${BUILDPACK_DIR}/bin"

  build

  util::print::success "** GO build completed **"
}

function build() {
  if [[ -f "${BUILDPACK_DIR}/run/main.go" ]]; then
    echo "- Building /run/main.go binary ..."
    GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o "${BUILDPACK_DIR}/bin/run" "${BUILDPACK_DIR}/run/main.go"
    echo "- Built /run/main.go binary built"

    # Create symlinks
    for name in detect build; do
      ln -sf "run" "${BUILDPACK_DIR}/bin/${name}"
    done
  else
    echo
    util::print::error "** GO Build Failed: No main.go file found in ${BUILDPACK_DIR}/run **"
  fi
}

main "$@"
