package pnpm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paketo-buildpacks/packit/v2"
)

func Build() packit.BuildFunc {
	return func(ctx packit.BuildContext) (packit.BuildResult, error) {
		// Install pnpm globally using npm
		cmd := exec.Command("npm", "install", "-g", "pnpm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to install pnpm: %w", err)
		}

		// provide pnpm as a layer
		layer, err := ctx.Layers.Get("pnpm")
		if err != nil {
			return packit.BuildResult{}, err
		}

		// make sure pnpm is available in the PATH
		layer.SharedEnv.Prepend("PATH", layer.Path, ":")

		return packit.BuildResult{
			Layers: []packit.Layer{layer},
		}, nil
	}
}
