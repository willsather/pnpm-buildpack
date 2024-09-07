package main

import (
	"github.com/paketo-buildpacks/packit/v2"
	pnpmstart "github.com/willsather/pnpm-buildpack/pnpm-start"
)

func main() {
	packit.Run(
		pnpmstart.Detect(),
		pnpmstart.Build(),
	)
}
