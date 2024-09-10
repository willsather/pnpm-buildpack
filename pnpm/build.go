package pnpm

import (
	"os"
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/draft"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate moq -out fakes/postal_service.go -pkg fakes . PostalService
type PostalService interface {
	Resolve(path string, id string, version string, stack string) (postal.Dependency, error)
	Deliver(dependency postal.Dependency, cnbPath string, layerPath string, platformPath string) error
}

func Build(
	postalService PostalService,
	logger scribe.Emitter,
	clock chronos.Clock,
) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		pnpmLayer, err := context.Layers.Get(Pnpm)
		if err != nil {
			return packit.BuildResult{}, err
		}

		planner := draft.NewPlanner()
		entry, _ := planner.Resolve(Pnpm, context.Plan.Entries, nil)
		version, ok := entry.Metadata["version"].(string)
		if !ok {
			version = "default"
		}

		dependency, err := postalService.Resolve(
			filepath.Join(context.CNBPath, "buildpack.toml"),
			entry.Name,
			version,
			context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		launch, build := planner.MergeLayerTypes("pnpm", context.Plan.Entries)

		var buildMetadata = packit.BuildMetadata{}
		var launchMetadata = packit.LaunchMetadata{}

		// check if we can reuse a previously built layer
		cachedSHA, ok := pnpmLayer.Metadata["checksum"].(string)
		if ok && postal.Checksum(dependency.Checksum).MatchString(cachedSHA) {
			logger.Process("Reusing cached layer %s", pnpmLayer.Path)
			logger.Break()

			pnpmLayer.Launch, pnpmLayer.Build, pnpmLayer.Cache = launch, build, build

			return packit.BuildResult{
				Layers: []packit.Layer{pnpmLayer},
				Build:  buildMetadata,
				Launch: launchMetadata,
			}, nil
		}

		logger.Process("Executing build process")

		pnpmLayer, err = pnpmLayer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		// FIXME: does this need to be set depending on the plan?
		//pnpmLayer.Launch, pnpmLayer.Build, pnpmLayer.Cache = launch, build, build
		pnpmLayer.Launch, pnpmLayer.Build, pnpmLayer.Cache = true, true, true

		logger.Subprocess("Installing pnpm v%s", dependency.Version)

		duration, err := clock.Measure(func() error {
			return postalService.Deliver(dependency, context.CNBPath, pnpmLayer.Path, context.Platform.Path)
		})
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		// set pnpm executable on layer path
		pnpmBinPath := filepath.Join(pnpmLayer.Path)
		pnpmLayer.SharedEnv.Append("PATH", pnpmBinPath, string(os.PathListSeparator))

		// set the layer metadata for future caching purposes
		pnpmLayer.Metadata = map[string]interface{}{
			"checksum": dependency.Checksum,
		}

		return packit.BuildResult{
			Layers: []packit.Layer{
				pnpmLayer,
			},
			Build:  buildMetadata,
			Launch: launchMetadata,
		}, nil
	}
}
