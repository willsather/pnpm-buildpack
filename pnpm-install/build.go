package pnpminstall

import (
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/libnodejs"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/draft"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate moq -out fakes/dependency_service.go -pkg fakes . DependencyService
type DependencyService interface {
	Install(path string, launch bool) error
}

func Build(dependencyService DependencyService, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("<<< Running PNPM Install Build")

		// step 1: get pnpm project path
		projectPath, err := libnodejs.FindProjectPath(context.WorkingDir)
		if err != nil {
			return packit.BuildResult{}, err
		}

		// store build and launch (if needed) layers to return
		var layers []packit.Layer

		// Step 2: check if installing deps for launch or build
		planner := draft.NewPlanner()
		launch, build := planner.MergeLayerTypes(NodeModules, context.Plan.Entries)

		// step 3a: install build dependencies, if required by plan
		if build {
			// create build-modules layer
			layer, err := context.Layers.Get("build-modules")
			if err != nil {
				return packit.BuildResult{}, err
			}

			layer, err = layer.Reset()
			if err != nil {
				return packit.BuildResult{}, err
			}

			// make new node_modules folder (if it doesn't exist)
			err = os.MkdirAll(filepath.Join(projectPath, "node_modules"), os.ModePerm)
			if err != nil {
				return packit.BuildResult{}, err
			}

			logger.Action("<> (Build) Installing Build Dependencies")
			err = dependencyService.Install(projectPath, false)

			// TODO: Add `install` cache with sha on the layer (and then check `shouldInstall` before actually installing

			// move node_modules folder onto layer
			err = fs.Move(filepath.Join(projectPath, "node_modules"), filepath.Join(layer.Path, "node_modules"))
			if err != nil {
				return packit.BuildResult{}, err
			}

			// symlink node_modules folder
			err = os.Symlink(filepath.Join(layer.Path, "node_modules"), filepath.Join(projectPath, "node_modules"))
			if err != nil {
				return packit.BuildResult{}, err
			}

			// attach `node_modules/.bin` to layer path
			nodeModulesBinPath := filepath.Join(layer.Path, "node_modules", ".bin")
			layer.LaunchEnv.Append("PATH", nodeModulesBinPath, string(os.PathListSeparator))
			layer.BuildEnv.Override("NODE_ENV", "development")

			logger.Detail("* Installed Build Dependencies")
			logger.Break()

			layer.Build = true
			layer.Cache = true

			// add build-modules layer to return result
			layers = append(layers, layer)
		}

		// step 3b: install launch dependencies, if required by plan
		if launch {
			// create launch-modules layer
			layer, err := context.Layers.Get("launch-modules")
			if err != nil {
				return packit.BuildResult{}, err
			}

			layer, err = layer.Reset()
			if err != nil {
				return packit.BuildResult{}, err
			}

			// make new node_modules folder (if it doesn't exist)
			err = os.MkdirAll(filepath.Join(projectPath, "node_modules"), os.ModePerm)
			if err != nil {
				return packit.BuildResult{}, err
			}

			logger.Action("<> (Build) Installing Launch Dependencies")
			err = dependencyService.Install(projectPath, false)

			// TODO: Add `install` cache with sha on the layer (and then check `shouldInstall` before actually installing

			// move node_modules folder onto layer
			err = fs.Move(filepath.Join(projectPath, "node_modules"), filepath.Join(layer.Path, "node_modules"))
			if err != nil {
				return packit.BuildResult{}, err
			}

			// symlink node_modules folder
			err = os.Symlink(filepath.Join(layer.Path, "node_modules"), filepath.Join(projectPath, "node_modules"))
			if err != nil {
				return packit.BuildResult{}, err
			}

			// attach `node_modules/.bin` to layer path
			nodeModulesBinPath := filepath.Join(layer.Path, "node_modules", ".bin")
			layer.LaunchEnv.Append("PATH", nodeModulesBinPath, string(os.PathListSeparator))
			layer.BuildEnv.Override("NODE_ENV", "production")

			logger.Detail("* Installed Build Dependencies")
			logger.Break()

			layer.Launch = true

			// add launch-modules layer to return result
			layers = append(layers, layer)
		}

		// step 3: return build-modules and launch-module layers (if present)
		return packit.BuildResult{
			Layers: layers,
		}, nil
	}
}
