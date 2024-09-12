package pnpm_test

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/willsather/pnpm-buildpack/pnpm"
)

func Test_DetectSuccessfully(t *testing.T) {
	var Expect = NewWithT(t).Expect
	var workingDirectory = t.TempDir()

	// setup mock logger
	var mockLogger = scribe.NewEmitter(bytes.NewBuffer(nil))
	var detect = pnpm.Detect(mockLogger)

	// when detect is called
	result, err := detect(packit.DetectContext{
		WorkingDir: workingDirectory,
	})

	// then detect returns a successful build plan
	Expect(err).NotTo(HaveOccurred())
	Expect(result.Plan).To(Equal(packit.BuildPlan{
		Provides: []packit.BuildPlanProvision{
			{
				Name: "pnpm",
			},
		},
	}))
}
