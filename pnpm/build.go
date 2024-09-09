package pnpm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(ctx packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("<<< Running PNPM Build")

		logger.Action("<> Installing PNPM Globally")

		// Install pnpm globally using npm
		cmd := exec.Command("npm", "install", "-g", "pnpm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to install pnpm: %w", err)
		}

		logger.Detail("* Providing PNPM Layer")

		// provide pnpm as a pnpmLayer
		pnpmLayer, err := ctx.Layers.Get(Pnpm)
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to retrieve pnpm layer: %w", err)
		}

		// make sure pnpm is available in the PATH
		pnpmLayer.SharedEnv.Prepend("PATH", pnpmLayer.Path, ":")

		logger.Detail("* Returning PNPM Build Result")

		return packit.BuildResult{
			Layers: []packit.Layer{pnpmLayer},
		}, nil
	}
}
