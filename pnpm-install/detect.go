package pnpminstall

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/libnodejs"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

type BuildPlanMetadata struct {
	Version       string `toml:"version"`
	VersionSource string `toml:"version-source"`
	Build         bool   `toml:"build"`
}

// CustomPackageJSON represents the structure of the package.json file
type CustomPackageJSON struct {
	PackageManager string `json:"packageManager"`
}

// getPnpmVersion reads the package.json file from the given working directory,
// parses the JSON, and extracts the pnpm version if present.
func getPnpmVersion(workingDirectory string) (string, error) {
	packageJSONPath := filepath.Join(workingDirectory, "package.json")

	// Read the package.json file
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return "", fmt.Errorf("error reading package.json: %w", err)
	}

	// Unmarshal the JSON data into a struct
	var packageJson CustomPackageJSON
	err = json.Unmarshal(data, &packageJson)
	if err != nil {
		return "", fmt.Errorf("error parsing package.json: %w", err)
	}

	// if the packageManager is pnpm and extract the version
	if strings.HasPrefix(packageJson.PackageManager, "pnpm@") {
		version := strings.TrimPrefix(packageJson.PackageManager, "pnpm@")
		return version, nil
	}

	return "", fmt.Errorf("pnpm version not found in package.json")
}

func Detect(logger scribe.Emitter) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		logger.Title("<<< Running PNPM Install Detect")

		// retrieve working directory
		projectPath, err := libnodejs.FindProjectPath(context.WorkingDir)
		if err != nil {
			return packit.DetectResult{}, err
		}

		// check if `pnpm-lock.yaml` is present
		exists, err := fs.Exists(filepath.Join(projectPath, "pnpm-lock.yaml"))
		if err != nil {
			return packit.DetectResult{}, err
		}

		if !exists {
			return packit.DetectResult{}, packit.Fail.WithMessage("no 'pnpm-lock.yaml' file found in the project path %s", projectPath)
		}

		// check if `package.json` is present
		pkg, err := libnodejs.ParsePackageJSON(projectPath)
		if err != nil {
			if os.IsNotExist(err) {
				return packit.DetectResult{}, packit.Fail.WithMessage("no 'package.json' found in project path %s", projectPath)
			}
			return packit.DetectResult{}, fmt.Errorf("failed to open package.json: %w", err)
		}

		// check if `package.json` has a `start` script present
		if !pkg.HasStartScript() {
			return packit.DetectResult{}, packit.Fail.WithMessage("'package.json' has been found but does not have a 'start' command")
		}

		nodeVersion := pkg.GetVersion()
		pnpmVersion, err := getPnpmVersion(projectPath)

		// build node requirement and metadata
		nodeRequirement := packit.BuildPlanRequirement{
			Name: Node,
			Metadata: BuildPlanMetadata{
				Build: true,
			},
		}

		if nodeVersion != "" {
			nodeRequirement.Metadata = BuildPlanMetadata{
				Version:       nodeVersion,
				VersionSource: "package.json",
				Build:         true,
			}
		}

		// build pnpm requirement and metadata
		pnpmRequirement := packit.BuildPlanRequirement{
			Name: Pnpm,
			Metadata: BuildPlanMetadata{
				Build: true,
			},
		}

		if pnpmVersion != "" {
			pnpmRequirement.Metadata = BuildPlanMetadata{
				Version:       pnpmVersion,
				VersionSource: "package.json",
				Build:         true,
			}
		}

		logger.Detail("* Return Build Plan that provides 'node_modules'")

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: NodeModules,
					},
				},
				Requires: []packit.BuildPlanRequirement{
					nodeRequirement,
					pnpmRequirement,
				},
			},
		}, nil
	}
}
