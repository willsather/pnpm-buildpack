package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var config struct {
	Buildpacks struct {
		BuildPlan string // used to read `plan.toml` into build plan
		Pnpm      string // used to install `pnpm` on layer
	}
}

func TestIntegration(t *testing.T) {
	var err error

	Expect := NewWithT(t).Expect

	// load in buildpack configuration into `config` struct
	file, err := os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(file).Decode(&config)
	Expect(err).NotTo(HaveOccurred())

	// get `pnpm` buildpack path
	root, err := filepath.Abs("./..")
	Expect(err).ToNot(HaveOccurred())

	buildpackStore := occam.NewBuildpackStore()

	// build and package pnpm buildpack
	config.Buildpacks.Pnpm, err = buildpackStore.Get.
		WithVersion("0.0.1").
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	// pull build plan buildpack
	config.Buildpacks.BuildPlan, err = buildpackStore.Get.
		Execute("github.com/paketo-community/build-plan")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	// run all integration test suites
	suite := spec.New("Integration", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Default", testDefault)
	suite.Run(t)
}
