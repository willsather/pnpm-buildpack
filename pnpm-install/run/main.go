package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/willsather/pnpm-buildpack/pnpm-install"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout)

	packit.Run(
		pnpminstall.Detect(logger),
		pnpminstall.Build(logger),
	)
}
