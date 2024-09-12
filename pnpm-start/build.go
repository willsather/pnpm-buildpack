package pnpmstart

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"os"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("<<< Running PNPM Start Build")

		// FIXME: support `prestart` and `poststart` commands like npm and yarn do
		//projectPath, err := libnodejs.FindProjectPath(context.WorkingDir)
		//if err != nil {
		//	return packit.BuildResult{}, err
		//}

		//logger.Action("<> Parsing Package Json")

		//pkg, err := libnodejs.ParsePackageJSON(projectPath)
		//if err != nil {
		//	return packit.BuildResult{}, err
		//}

		//command := "pnpm"
		//arg := fmt.Sprintf("pnpm start")

		//if pkg.Scripts.Start != "" {
		//	command = "pnpm"
		//	arg = fmt.Sprintf("pnpm start")
		//}

		//if pkg.Scripts.PreStart != "" {
		//	command = "bash"
		//	arg = fmt.Sprintf("%s && %s", pkg.Scripts.PreStart, arg)
		//}
		//
		//if pkg.Scripts.PostStart != "" {
		//	command = "bash"
		//	arg = fmt.Sprintf("%s && %s", arg, pkg.Scripts.PostStart)
		//}

		// Ideally we would like the lifecycle to support setting a custom working
		// directory to run the launch process.  Until that happens we will cd in.
		//if projectPath != context.WorkingDir {
		//	command = "bash"
		//	arg = fmt.Sprintf("cd %s && %s", projectPath, arg)
		//}
		//
		//args := []string{arg}
		//if command == "bash" {
		//	args = []string{"-c", arg}
		//}

		processes := []packit.Process{
			{
				Type:    "web",
				Command: "pnpm",
				Args:    []string{"start"},
				Default: true,
				Direct:  true,
			},
		}

		logger.LaunchProcesses(processes)
		logger.Detail("* PNPM Start Build returning Processes")
		logger.Detail("* PATH:", os.Getenv("PATH"))

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
