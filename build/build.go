package build

import (
	"os/exec"

	"github.com/paketo-buildpacks/packit"
)

func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		cmd := exec.Command("pnpm", "install")
		cmd.Stdout = context.Stdout
		cmd.Stderr = context.Stderr

		if err := cmd.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		// Create the launch process
		process := packit.Process{
			Type:    "web",
			Command: "pnpm",
			Args:    []string{"start"},
			Direct:  true,
		}

		return packit.BuildResult{
			Layers: []packit.Layer{},
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{process},
			},
		}, nil
	}
}
