package pnpmstart_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/willsather/pnpm-buildpack/pnpm-start"
)

var build = pnpmstart.Build(scribe.NewEmitter(bytes.NewBuffer(nil)))

func Test_BuildSuccessfully(t *testing.T) {
	var Expect = NewWithT(t).Expect

	// setup working directory
	var workingDirectory, err = os.MkdirTemp("", "working-directory")
	Expect(err).NotTo(HaveOccurred())

	// setup project directory and files
	Expect(os.Mkdir(filepath.Join(workingDirectory, "project"), os.ModePerm)).To(Succeed())

	// WHEN: build is called
	result, err := build(packit.BuildContext{
		WorkingDir: filepath.Join(workingDirectory, "project"),
	})

	// THEN: no error and launch plan includes processes
	Expect(err).NotTo(HaveOccurred())
	Expect(result.Launch).To(Equal(packit.LaunchMetadata{
		Processes: []packit.Process{{
			Type:    "web",
			Command: "pnpm",
			Args:    []string{"start"},
			Default: true,
			Direct:  true,
		}},
	}))

	// cleanup
	Expect(os.RemoveAll(workingDirectory)).To(Succeed())
}
