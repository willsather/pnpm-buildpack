package pnpm_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/willsather/pnpm-buildpack/pnpm"
	"github.com/willsather/pnpm-buildpack/pnpm/fakes"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

// setup mock logger and clock
var mockLogger = scribe.NewEmitter(bytes.NewBuffer(nil))

func Test_BuildSuccessfully(t *testing.T) {
	// GIVEN
	var Expect = NewWithT(t).Expect
	var workingDirectory = t.TempDir()

	// setup mock postal service
	mockPostalService := &fakes.PostalServiceMock{
		ResolveFunc: func(path string, id string, version string, stack string) (postal.Dependency, error) {
			return postal.Dependency{
				ID:       "pnpm",
				Name:     "pnpm",
				Checksum: "sha256:fake-pnpm-exe-checksum",
				Stacks:   []string{"*"},
				URI:      "fake-download-pnpm-exe-url",
				Version:  "fake-pnpm-version",
			}, nil
		},
		DeliverFunc: func(dependency postal.Dependency, cnbPath string, layerPath string, platformPath string) error {
			return nil
		},
	}

	// setup mock build context directories
	cnbDirectory, err := os.MkdirTemp("", "cnb")
	Expect(err).NotTo(HaveOccurred())

	layersDirectory, err := os.MkdirTemp("", "layers")
	Expect(err).NotTo(HaveOccurred())

	// setup build function with dependencies
	var build = pnpm.Build(mockPostalService, mockLogger, chronos.DefaultClock)

	// WHEN: build is called (plan includes what was generated in detect)
	result, err := build(packit.BuildContext{
		WorkingDir: workingDirectory,
		CNBPath:    cnbDirectory,
		Stack:      "*",
		BuildpackInfo: packit.Info{
			Name:    "PNPM Buildpack",
			Version: "0.0.1",
		},
		Plan: packit.BuildpackPlan{
			Entries: []packit.BuildpackPlanEntry{
				{
					Name: "pnpm",
				},
			},
		},
		Platform: packit.Platform{Path: "platform"},
		Layers:   packit.Layers{Path: layersDirectory},
	})

	// then build returns a successful build plan
	Expect(err).NotTo(HaveOccurred())

	// then build plan contains valid pnpm layer
	Expect(result.Layers).To(HaveLen(1))
	layer := result.Layers[0]

	Expect(layer.Name).To(Equal("pnpm"))
	Expect(layer.Path).To(Equal(filepath.Join(layersDirectory, "pnpm")))
	Expect(layer.Metadata).To(Equal(map[string]interface{}{
		"checksum": "sha256:fake-pnpm-exe-checksum",
	}))

	// then postal service resolve is called (once, with parameters)
	Expect(len(mockPostalService.ResolveCalls())).To(Equal(1))
	Expect(mockPostalService.ResolveCalls()[0].ID).To(Equal("pnpm"))
	Expect(mockPostalService.ResolveCalls()[0].Stack).To(Equal("*"))
	Expect(mockPostalService.ResolveCalls()[0].Path).To(Equal(filepath.Join(cnbDirectory, "buildpack.toml")))

	// then postal service deliver is called (once, with parameters)
	Expect(len(mockPostalService.DeliverCalls())).To(Equal(1))
	Expect(mockPostalService.DeliverCalls()[0].CnbPath).To(Equal(cnbDirectory))
	Expect(mockPostalService.DeliverCalls()[0].LayerPath).To(Equal(filepath.Join(layersDirectory, "pnpm")))
	Expect(mockPostalService.DeliverCalls()[0].PlatformPath).To(Equal("platform"))
	Expect(mockPostalService.DeliverCalls()[0].Dependency).To(Equal(postal.Dependency{
		ID:       "pnpm",
		Name:     "pnpm",
		Checksum: "sha256:fake-pnpm-exe-checksum",
		Stacks:   []string{"*"},
		URI:      "fake-download-pnpm-exe-url",
		Version:  "fake-pnpm-version",
	}))

	// cleanup
	Expect(os.RemoveAll(layersDirectory)).To(Succeed())
	Expect(os.RemoveAll(cnbDirectory)).To(Succeed())
	Expect(os.RemoveAll(workingDirectory)).To(Succeed())
}

func Test_BuildFailsDependencyCannotResolve(t *testing.T) {
	// GIVEN
	var Expect = NewWithT(t).Expect
	var workingDirectory = t.TempDir()

	// setup mock postal service that returns resolve error
	mockPostalService := &fakes.PostalServiceMock{
		ResolveFunc: func(path string, id string, version string, stack string) (postal.Dependency, error) {
			return postal.Dependency{}, errors.New("Postal.Resolve cannot resolve dependency")
		},
	}

	// setup mock build context directories
	cnbDirectory, err := os.MkdirTemp("", "cnb")
	Expect(err).NotTo(HaveOccurred())

	layersDirectory, err := os.MkdirTemp("", "layers")
	Expect(err).NotTo(HaveOccurred())

	// setup build function with dependencies
	var build = pnpm.Build(mockPostalService, mockLogger, chronos.DefaultClock)

	// WHEN: build is called (plan includes what was generated in detect)
	_, err = build(packit.BuildContext{
		WorkingDir: workingDirectory,
		CNBPath:    cnbDirectory,
		Stack:      "*",
		BuildpackInfo: packit.Info{
			Name:    "PNPM Buildpack",
			Version: "0.0.1",
		},
		Plan: packit.BuildpackPlan{
			Entries: []packit.BuildpackPlanEntry{
				{
					Name: "pnpm",
				},
			},
		},
		Platform: packit.Platform{Path: "platform"},
		Layers:   packit.Layers{Path: layersDirectory},
	})

	// then build returns a successful build plan
	Expect(err).To(HaveOccurred())
	Expect(err).To(MatchError("Postal.Resolve cannot resolve dependency"))

	// cleanup
	Expect(os.RemoveAll(layersDirectory)).To(Succeed())
	Expect(os.RemoveAll(cnbDirectory)).To(Succeed())
	Expect(os.RemoveAll(workingDirectory)).To(Succeed())
}

// TODO: Test more: that cache layer doesn't call Deliver
// TODO: More failure cases?
