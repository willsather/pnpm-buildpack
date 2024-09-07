package pnpm

import (
	"fmt"
	"github.com/paketo-buildpacks/packit/v2"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		fmt.Println("<<< Running PNPM Detect >>>")

		fmt.Println("<<< Returning Build Plan that provides PNPM >>>")
		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: Pnpm},
				},
			},
		}, nil
	}
}
