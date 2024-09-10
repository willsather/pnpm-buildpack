package pnpminstall

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("<<< Running PNPM Install Build")

		// FIXME: are we okay with avoiding this now?
		// Maybe we do this on `debug` logs?

		// step 0: verify npm and pnpm were installed
		vNPM := exec.Command("npm", "-v")
		vNPM.Stdout = os.Stdout
		vNPM.Stderr = os.Stderr

		if err := vNPM.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		vPNPM := exec.Command("pnpm", "-v")
		vPNPM.Stdout = os.Stdout
		vPNPM.Stderr = os.Stderr

		if err := vPNPM.Run(); err != nil {
			return packit.BuildResult{}, err
		}

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
		logger.Action("<> (Build) Installing All Dependencies")

		// get layer for the dependencies required for build
		layer, err := context.Layers.Get("build-modules")
		if err != nil {
			return packit.BuildResult{}, err
		}

		// maybe this needs to be reset?
		layer, err = layer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		// make new node_modules folder
		err = os.Mkdir(filepath.Join(layer.Path, "node_modules"), os.ModePerm)
		if err != nil {
			return packit.BuildResult{}, err
		}

		installCmd := exec.Command("pnpm", "install", "--frozen-lockfile")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		if err := installCmd.Run(); err != nil {
			return packit.BuildResult{}, err
		}

		// FIXME: are these required?
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
