package main

import (
	"github.com/paketo-buildpacks/packit/v2"
	pnpmbuildpack "github.com/willsather/pnpm-buildpack/pnpm-install"
)

func main() {
	packit.Run(
		pnpmbuildpack.Detect(),
		pnpmbuildpack.Build(),
	)
}
