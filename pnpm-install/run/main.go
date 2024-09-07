package main

import (
	"github.com/paketo-buildpacks/packit/v2"
	pnpminstall "github.com/willsather/pnpm-buildpack/pnpm-install"
)

func main() {
	packit.Run(
		pnpminstall.Detect(),
		pnpminstall.Build(),
	)
}
