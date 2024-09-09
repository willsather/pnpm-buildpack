package pnpmstart

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/libnodejs"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Detect(logger scribe.Emitter) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		logger.Title("<<< Running PNPM Start Detect")

		// find working directory
		projectPath, err := libnodejs.FindProjectPath(context.WorkingDir)
		if err != nil {
			return packit.DetectResult{}, err
		}

		// find pnpm-lock.yaml
		exists, err := fs.Exists(filepath.Join(projectPath, "pnpm-lock.yaml"))
		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("failed to stat pnpm-lock.yaml: %w", err)
		}

		if !exists {
			return packit.DetectResult{}, packit.Fail.WithMessage("no 'pnpm-lock.yaml' found in the project path %s", projectPath)
		}

		// find package.json
		pkg, err := libnodejs.ParsePackageJSON(projectPath)
		if err != nil {
			if os.IsNotExist(err) {
				return packit.DetectResult{}, packit.Fail.WithMessage("no 'package.json' found in project path %s", projectPath)
			}
			return packit.DetectResult{}, fmt.Errorf("failed to open package.json: %w", err)
		}

		// check if start script exists
		if !pkg.HasStartScript() {
			return packit.DetectResult{}, packit.Fail.WithMessage("'package.json' does not have a 'start' script")
		}

		logger.Detail("* PNPM Start Returning Build Plan with ONLY requirements")

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{
					{
						Name: Node,
						Metadata: map[string]interface{}{
							"launch": true,
						},
					},
					{
						Name: Pnpm,
						Metadata: map[string]interface{}{
							"launch": true,
						},
					},
					{
						Name: NodeModules,
						Metadata: map[string]interface{}{
							"launch": true,
						},
					},
				},
			},
		}, nil
	}
}
