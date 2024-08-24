package main

import (
	"github.com/paketo-buildpacks/packit"
	pnpmbuildpack "github.com/willsather/pnpm-buildpack"
)

func main() {
	packit.Run(
		pnpmbuildpack.Detect(),
		pnpmbuildpack.Build(),
	)
}
