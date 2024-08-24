#!/usr/bin/env bash

set -eu
set -o pipefail

# Directory locations
readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACKDIR="$(cd "${PROGDIR}/.." && pwd)"

# import utils
source "${PROGDIR}/utils/print.sh"

function main() {
  util::print::title "** GO build **"

  mkdir -p "${BUILDPACKDIR}/bin"

  build::run

  util::print::success "** GO build completed **"
}

function build::run() {
  if [[ -f "${BUILDPACKDIR}/run/main.go" ]]; then
    echo "- Building /run/main.go binary ..."
    GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o "${BUILDPACKDIR}/bin/run" "${BUILDPACKDIR}/run/main.go"
    echo "- Built /run/main.go binary built"

    # Create symlinks
    for name in detect build; do
      ln -sf "run" "${BUILDPACKDIR}/bin/${name}"
    done
  else
    echo
    util::print::error "** GO Build Failed: No main.go file found in ${BUILDPACKDIR}/run **"
  fi
}

function build::executable() {
  if [[ -d "${BUILDPACKDIR}/cmd" ]]; then
    for src in "${BUILDPACKDIR}/cmd/"*; do
      local name
      name="$(basename "${src}")"

      if [[ -f "${src}/main.go" ]]; then
        echo "Building ${name}..."
        GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o "${BUILDPACKDIR}/bin/${name}" "${src}/main.go"
        echo "${name} built successfully."
      else
        echo "Skipping ${name}, no main.go file."
      fi
    done
  fi
}

main "$@"
