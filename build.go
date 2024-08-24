package pnpm_buildpack

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paketo-buildpacks/packit"
)

func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		fmt.Println("<<< Starting Build Process >>>")

		// Step 1: Install dependencies
		installCmd := exec.Command("pnpm", "install")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		fmt.Println("<<< Installing Dependencies >>>")
		if err := installCmd.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		fmt.Println("<<< Installed Dependencies >>>")

		// Step 2: Optionally build the project
		buildCmd := exec.Command("pnpm", "build")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		fmt.Println("<<< Building Application >>>")

		if err := buildCmd.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		fmt.Println("<<< Built Application >>>")

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
