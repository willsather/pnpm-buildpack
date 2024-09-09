package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/willsather/pnpm-buildpack/pnpm-start"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout)

	packit.Run(
		pnpmstart.Detect(logger),
		pnpmstart.Build(logger),
	)
}
