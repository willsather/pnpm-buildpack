package pnpmstart

import (
	"fmt"

	"github.com/paketo-buildpacks/libnodejs"
	"github.com/paketo-buildpacks/packit/v2"
)

func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		fmt.Println("<<< Running PNPM Start Build >>>")

		projectPath, err := libnodejs.FindProjectPath(context.WorkingDir)
		if err != nil {
			return packit.BuildResult{}, err
		}

		fmt.Println("<<< Parsing Package Json >>>")

		pkg, err := libnodejs.ParsePackageJSON(projectPath)
		if err != nil {
			return packit.BuildResult{}, err
		}

		command := "pnpm"
		arg := fmt.Sprintf("pnpm start")

		if pkg.Scripts.Start != "" {
			command = "pnpm"
			arg = fmt.Sprintf("pnpm start")
		}

		if pkg.Scripts.PreStart != "" {
			command = "bash"
			arg = fmt.Sprintf("%s && %s", pkg.Scripts.PreStart, arg)
		}

		if pkg.Scripts.PostStart != "" {
			command = "bash"
			arg = fmt.Sprintf("%s && %s", arg, pkg.Scripts.PostStart)
		}

		// Ideally we would like the lifecycle to support setting a custom working
		// directory to run the launch process.  Until that happens we will cd in.
		if projectPath != context.WorkingDir {
			command = "bash"
			arg = fmt.Sprintf("cd %s && %s", projectPath, arg)
		}

		args := []string{arg}
		if command == "bash" {
			args = []string{"-c", arg}
		}

		// TODO: add more logging using the `packit` scribe (instead of `fmt`)
		// can call `logger.LaunchProcesses(processes)` here

		fmt.Println("<<< PNPM Start Build returning Processes >>>")

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type:    "web",
						Command: command,
						Args:    args,
						Default: true,
						Direct:  true,
					},
				},
			},
		}, nil
	}
}
