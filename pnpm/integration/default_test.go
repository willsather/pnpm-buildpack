package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDefault(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually
		err        error

		pack   occam.Pack
		docker occam.Docker

		pullPolicy = "always"

		image occam.Image

		name   string
		source string

		container occam.Container
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose()
		docker = occam.NewDocker()

		name, err = occam.RandomName()
		Expect(err).NotTo(HaveOccurred())

		source, err = occam.Source(filepath.Join("collateral", "simple-app"))
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())

		Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
		Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
		Expect(os.RemoveAll(source)).To(Succeed())
	})

	it("builds basic image with pnpm", func() {
		// GIVEN: occam and docker setup is complete
		var (
			logs fmt.Stringer
			err  error
		)

		// WHEN: buildpack is used to package application
		image, logs, err = pack.WithNoColor().Build.
			WithBuildpacks(
				config.Buildpacks.Pnpm,
				config.Buildpacks.BuildPlan,
			).
			WithPullPolicy(pullPolicy).
			Execute(name, source)

		// THEN: no error occurred
		Expect(err).ToNot(HaveOccurred(), logs.String)

		// THEN: pnpm is present
		container, err = docker.Container.Run.WithCommand("command -v pnpm").Execute(image.ID)
		Expect(err).NotTo(HaveOccurred())

		// THEN: pnpm is found in logs
		Eventually(func() string {
			cLogs, err := docker.Container.Logs.Execute(container.ID)
			Expect(err).NotTo(HaveOccurred())
			return cLogs.String()
		}).Should(ContainSubstring("pnpm"))
	})

}
