package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
)

func Test_SimpleSampleApplicationPnpmInstall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		image     occam.Image
		container occam.Container

		buildpackInfo struct {
			Buildpack struct {
				ID   string
				Name string
			}
		}

		err error
	)

	// before
	pack := occam.NewPack()
	docker := occam.NewDocker()

	name, err := occam.RandomName()
	Expect(err).NotTo(HaveOccurred())

	// setup other buildpacks
	file, err := os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(file).Decode(&buildpackInfo)
	Expect(err).NotTo(HaveOccurred())

	// create occam store to pull all necessary buildpacks
	buildpackStore := occam.NewBuildpackStore()

	nodeURI, err := buildpackStore.Get.Execute("github.com/paketo-buildpacks/node-engine")
	Expect(err).ToNot(HaveOccurred())

	// get `pnpm` buildpack path
	pnpmDirectory, err := filepath.Abs("./../../pnpm")
	Expect(err).ToNot(HaveOccurred())

	// get `pnpm` buildpack
	pnpmURI, err := buildpackStore.Get.WithVersion("0.0.1").Execute(pnpmDirectory)
	Expect(err).ToNot(HaveOccurred())

	// source the sample project
	source, err := occam.Source(filepath.Join("collateral", "sample-app"))
	Expect(err).NotTo(HaveOccurred())

	// when build is called
	image, _, err = pack.Build.
		WithBuildpacks(
			nodeURI,
			pnpmURI,
		).
		WithPullPolicy("always").
		Execute(name, source)

	Expect(err).NotTo(HaveOccurred())

	// check the contents of the node modules
	container, err = docker.Container.Run.
		WithCommand(fmt.Sprintf("ls -alR /layers/%s/launch-modules/node_modules",
			strings.ReplaceAll(buildpackInfo.Buildpack.ID, "/", "_"))).
		Execute(image.ID)

	Expect(err).NotTo(HaveOccurred())

	Eventually(func() string {
		cLogs, err := docker.Container.Logs.Execute(container.ID)
		Expect(err).NotTo(HaveOccurred())
		return cLogs.String()
	}).Should(ContainSubstring("is-even"))

	// cleanup
	Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
	Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
	Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
	Expect(os.RemoveAll(source)).To(Succeed())
}
