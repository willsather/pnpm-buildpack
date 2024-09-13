package pnpminstall_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/willsather/pnpm-buildpack/pnpm-install"
	"github.com/willsather/pnpm-buildpack/pnpm-install/fakes"
)

func Test_BuildSuccessfully(t *testing.T) {
	var Expect = NewWithT(t).Expect
	var workingDirectory = t.TempDir()

	mockDependencyService := &fakes.DependencyServiceMock{
		InstallFunc: func(cnbPath string) error {
			return nil
		},
	}

	// setup mock dependency service and mock logger
	var mockLogger = scribe.NewEmitter(bytes.NewBuffer(nil))
	var build = pnpminstall.Build(mockDependencyService, mockLogger)

	// setup mock pnpm project
	Expect(os.Mkdir(filepath.Join(workingDirectory, "project"), os.ModePerm)).To(Succeed())

	// setup mock build context directories
	cnbDirectory, err := os.MkdirTemp("", "cnb")
	Expect(err).NotTo(HaveOccurred())

	layersDirectory, err := os.MkdirTemp("", "layers")
	Expect(err).NotTo(HaveOccurred())

	// WHEN: pnpm install build is called
	result, err := build(packit.BuildContext{
		BuildpackInfo: packit.Info{
			Name:    "PNPM Buildpack",
			Version: "0.0.1",
		},
		WorkingDir: filepath.Join(workingDirectory, "project"),
		CNBPath:    cnbDirectory,
		Layers: packit.Layers{
			Path: layersDirectory,
		},
		Plan: packit.BuildpackPlan{
			Entries: []packit.BuildpackPlanEntry{
				{
					Name: "node_modules",
					Metadata: map[string]interface{}{
						"build": true,
					},
				},
			},
		},
		Stack: "*",
		Platform: packit.Platform{
			Path: "some-path",
		},
	})

	// then detect returns a successful build plan
	Expect(err).NotTo(HaveOccurred())

	// then the build-modules layer is created correctly
	Expect(len(result.Layers)).To(Equal(1))

	layer := result.Layers[0]
	Expect(layer.Name).To(Equal("build-modules"))
	Expect(layer.Path).To(Equal(filepath.Join(layersDirectory, "build-modules")))

	// then node_modules should be added to class path
	Expect(layer.LaunchEnv).To(Equal(packit.Environment{
		"PATH.append": filepath.Join(layersDirectory, "build-modules", "node_modules", ".bin"),
		"PATH.delim":  ":",
	}))

	// TODO: should assertions be made on `build` and `cache`?
	Expect(layer.Launch).To(BeTrue())

	// then dependency service `install` is called (once, with parameters)
	Expect(len(mockDependencyService.InstallCalls())).To(Equal(1))
	Expect(mockDependencyService.InstallCalls()[0].Path).To(Equal(filepath.Join(workingDirectory, "project")))

	// cleanup
	Expect(os.RemoveAll(layersDirectory)).To(Succeed())
	Expect(os.RemoveAll(workingDirectory)).To(Succeed())
	Expect(os.RemoveAll(cnbDirectory)).To(Succeed())
}
