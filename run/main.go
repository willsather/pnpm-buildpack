package main

import (
	"github.com/paketo-buildpacks/packit/v2"
	pnpmbuildpack "github.com/willsather/pnpm-buildpack"
)

func main() {
	packit.Run(
		pnpmbuildpack.Detect(),
		pnpmbuildpack.Build(),
	)
}
