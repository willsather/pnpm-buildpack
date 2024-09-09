package pnpminstall

import (
	"os"
	"os/exec"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("<<< Running PNPM Install Build")

		// FIXME: are we okay with avoiding this now?
		// Maybe we do this on `debug` logs?
		// step 0: verify npm and pnpm were installed
		vNPM := exec.Command("npm", "-v")
		vNPM.Stdout = os.Stdout
		vNPM.Stderr = os.Stderr

		if err := vNPM.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		vPNPM := exec.Command("pnpm", "-v")
		vPNPM.Stdout = os.Stdout
		vPNPM.Stderr = os.Stderr

		if err := vPNPM.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		// Step 1: Install dependencies
		logger.Action("<> Installing Dependencies")

		installCmd := exec.Command("pnpm", "install")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		if err := installCmd.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		logger.Detail("* Installed Dependencies")

		// Return the build result with the cached layer
		return packit.BuildResult{}, nil

		// TODO: do I need to use build/launch here? How does this work
		//// Step 2: Optionally build the project
		//buildCmd := exec.Command("pnpm", "build")
		//buildCmd.Stdout = os.Stdout
		//buildCmd.Stderr = os.Stderr
		//
		//fmt.Println("<<< Building Application >>>")
		//
		//if err := buildCmd.Run(); err != nil {
		//	return packit.BuildResult{}, err
		//}
		//
		//fmt.Println("<<< Built Application >>>")
		//
		//// Create the launch process
		//process := packit.Process{
		//	Type:    "web",
		//	Command: "pnpm",
		//	Args:    []string{"start"},
		//	Direct:  true,
		//}
		//
		//return packit.BuildResult{
		//	Layers: []packit.Layer{},
		//	Launch: packit.LaunchMetadata{
		//		Processes: []packit.Process{process},
		//	},
		//}, nil
	}
}
