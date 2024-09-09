package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/willsather/pnpm-buildpack/pnpm"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout)

	packit.Run(
		pnpm.Detect(logger),
		pnpm.Build(logger),
	)
}
