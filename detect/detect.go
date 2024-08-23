package detect

import (
	"os"

	"github.com/paketo-buildpacks/packit"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		if _, err := os.Stat("pnpm-lock.yaml"); os.IsNotExist(err) {
			return packit.DetectResult{}, packit.Fail.WithMessage("pnpm-lock.yaml not found")
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{{Name: "node_modules"}},
				Requires: []packit.BuildPlanRequirement{
					{Name: "node", Metadata: map[string]interface{}{"version": "16.*"}},
					{Name: "pnpm"},
				},
			},
		}, nil
	}
}
