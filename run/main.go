package main

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/willsather/pnpm-buildpack/build"
	"github.com/willsather/pnpm-buildpack/detect"
)

func main() {
	packit.Run(
		detect.Detect(),
		build.Build(),
	)
}
