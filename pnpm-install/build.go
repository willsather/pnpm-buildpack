package pnpminstall

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paketo-buildpacks/packit/v2"
)

func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		fmt.Println("<<< Starting PNPM Install Build Process >>>")

		// Step 1: Install dependencies
		installCmd := exec.Command("pnpm", "install")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		fmt.Println("<<< Installing Dependencies >>>")
		if err := installCmd.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		fmt.Println("<<< Installed Dependencies >>>")

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
