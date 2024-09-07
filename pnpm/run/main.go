package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	pnpmbuildpack "github.com/willsather/pnpm-buildpack/pnpm"
)

func main() {
	// Use os.Args to determine whether we're running detect or build
	command := os.Args[0]

	switch command {
	case "detect":
		packit.Detect(pnpmbuildpack.Detect())
	case "build":
		packit.Build(pnpmbuildpack.Build())
	default:
		os.Exit(1)
	}
}
