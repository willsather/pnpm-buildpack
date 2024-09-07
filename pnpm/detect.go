package pnpm

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

func Detect() packit.DetectFunc {
	return func(ctx packit.DetectContext) (packit.DetectResult, error) {
		// Check for the presence of the pnpm-lock.yaml file
		if exists, err := fs.Exists("pnpm-lock.yaml"); err != nil {
			return packit.DetectResult{}, err
		} else if exists {
			return packit.DetectResult{
				Plan: packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{{Name: "pnpm"}},
				},
			}, nil
		}
		return packit.DetectResult{}, nil
	}
}
