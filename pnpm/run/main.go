package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/willsather/pnpm-buildpack/pnpm"
)

func main() {
	// create logger
	logger := scribe.NewEmitter(os.Stdout)

	// create postal service (used to install dependencies in `build.go`)
	transport := cargo.NewTransport()
	postalService := postal.NewService(transport)

	packit.Run(
		pnpm.Detect(logger),
		pnpm.Build(postalService, logger, chronos.DefaultClock),
	)
}
