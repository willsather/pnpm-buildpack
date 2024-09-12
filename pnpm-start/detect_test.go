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

var detect = pnpmstart.Detect(scribe.NewEmitter(bytes.NewBuffer(nil)))

func Test_DetectSuccessfully(t *testing.T) {
	var Expect = NewWithT(t).Expect

	// setup working directory
	var workingDirectory, err = os.MkdirTemp("", "working-directory")
	Expect(err).NotTo(HaveOccurred())

	// setup project directory and files
	Expect(os.Mkdir(filepath.Join(workingDirectory, "project"), os.ModePerm)).To(Succeed())

	Expect(os.WriteFile(filepath.Join(workingDirectory, "project", "pnpm-lock.yaml"), nil, 0600)).To(Succeed())

	Expect(os.WriteFile(filepath.Join(workingDirectory, "project", "package.json"), []byte(`{
				"scripts": {
					"start": "node server.js"
				}
			}`), 0600)).To(Succeed())

	// WHEN: detect is called
	result, err := detect(packit.DetectContext{
		WorkingDir: filepath.Join(workingDirectory, "project"),
	})

	// THEN:
	Expect(err).NotTo(HaveOccurred())
	Expect(result.Plan).To(Equal(packit.BuildPlan{
		Requires: []packit.BuildPlanRequirement{
			{
				Name: "node",
				Metadata: map[string]interface{}{
					"launch": true,
				},
			},
			{
				Name: "pnpm",
				Metadata: map[string]interface{}{
					"launch": true,
				},
			},
			{
				Name: "node_modules",
				Metadata: map[string]interface{}{
					"launch": true,
				},
			},
		}}))

	// cleanup
	Expect(os.RemoveAll(workingDirectory)).To(Succeed())
}

func Test_DetectFailsWithoutPackageLock(t *testing.T) {
	var Expect = NewWithT(t).Expect

	// setup working directory
	var workingDirectory, err = os.MkdirTemp("", "working-directory")
	Expect(err).NotTo(HaveOccurred())

	// setup project directory
	Expect(os.Mkdir(filepath.Join(workingDirectory, "project"), os.ModePerm)).To(Succeed())

	// WHEN: detect is called
	_, err = detect(packit.DetectContext{
		WorkingDir: filepath.Join(workingDirectory, "project"),
	})

	// THEN: a failure arises
	Expect(err).To(MatchError(packit.Fail.WithMessage("no 'pnpm-lock.yaml' found in the project path %s", filepath.Join(workingDirectory, "project"))))

	// cleanup
	Expect(os.RemoveAll(workingDirectory)).To(Succeed())
}

func Test_DetectFailsWithoutPackageJson(t *testing.T) {
	var Expect = NewWithT(t).Expect

	// setup working directory
	var workingDirectory, err = os.MkdirTemp("", "working-directory")
	Expect(err).NotTo(HaveOccurred())

	// setup project directory and pnpm-lock
	Expect(os.Mkdir(filepath.Join(workingDirectory, "project"), os.ModePerm)).To(Succeed())

	Expect(os.WriteFile(filepath.Join(workingDirectory, "project", "pnpm-lock.yaml"), nil, 0600)).To(Succeed())

	// WHEN: detect is called
	_, err = detect(packit.DetectContext{
		WorkingDir: filepath.Join(workingDirectory, "project"),
	})

	// THEN: a failure arises
	Expect(err).To(MatchError(packit.Fail.WithMessage("no 'package.json' found in project path %s", filepath.Join(workingDirectory, "project"))))

	// cleanup
	Expect(os.RemoveAll(workingDirectory)).To(Succeed())
}
