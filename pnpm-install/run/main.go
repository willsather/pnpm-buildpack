package main

import (
	"os"
	"os/exec"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	pnpminstall "github.com/willsather/pnpm-buildpack/pnpm-install"
)

type DependencyService struct{}

func (s *DependencyService) Install(path string) error {
	installCmd := exec.Command("pnpm", "install", "--frozen-lockfile")
	installCmd.Dir = path
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	return installCmd.Run()
}

func main() {
	logger := scribe.NewEmitter(os.Stdout)
	dependencyService := &DependencyService{}

	packit.Run(
		pnpminstall.Detect(logger),
		pnpminstall.Build(dependencyService, logger),
	)
}
