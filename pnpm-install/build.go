package pnpminstall

import (
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate moq -out fakes/dependency_service.go -pkg fakes . DependencyService
type DependencyService interface {
	Install(path string) error
}

func Build(dependencyService DependencyService, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("<<< Running PNPM Install Build")

		// step XX: get pnpm project path
		//projectPath, err := libnodejs.FindProjectPath(context.WorkingDir)
		//if err != nil {
		//	return packit.BuildResult{}, err
		//}

		// store layers to return at end
		var layers []packit.Layer

		// Step XX: check if installing deps for launch or build
		//planner := draft.NewPlanner()

		// TODO: Add `launch` and `build` back to this first variable
		// launch, build := planner.MergeLayerTypes(NodeModules, context.Plan.Entries)

		//if build {

		// create layer for the dependencies required for build
		layer, err := context.Layers.Get("build-modules")
		if err != nil {
			return packit.BuildResult{}, err
		}

		// TODO: test if this is even needed ???
		layer, err = layer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		// make new node_modules folder (TODO: test context.workingDir vs layer.Path)
		err = os.Mkdir(filepath.Join(context.WorkingDir, "node_modules"), os.ModePerm)
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.Action("<> (Build) Installing All Dependencies")
		err = dependencyService.Install(filepath.Join(context.WorkingDir))

		// TODO: Add `install` cache with sha on the layer (and then check `shouldInstall` before actually installing

		// FIXME: is these symlinks / folder copy required?
		// move node_modules folder
		//err = fs.Move(filepath.Join(projectPath, "node_modules"), filepath.Join(layer.Path, "node_modules"))
		//if err != nil {
		//	return packit.BuildResult{}, err
		//}
		//
		//// symlink node_modules folder
		//err = os.Symlink(filepath.Join(layer.Path, "node_modules"), filepath.Join(projectPath, "node_modules"))
		//if err != nil {
		//	return packit.BuildResult{}, err
		//}

		logger.Detail("* Installed Dependencies")
		logger.Break()

		layer.Launch = true
		layers = append(layers, layer)
		//}

		// FIXME: don't install dev deps if launch is set
		// this shouldn't break launches, but rather just increase image size, etc
		//if launch {
		//	logger.Action("<> (Launch) Installing Production Dependencies")
		//	layer, err := context.Layers.Get("launch-modules")
		//
		//}

		return packit.BuildResult{
			Layers: layers,
		}, nil
	}
}
