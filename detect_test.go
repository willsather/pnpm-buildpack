package pnpm_buildpack

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit"
)

var detect = Detect()

func Test_DetectSuccessfully(t *testing.T) {
	var Expect = NewWithT(t).Expect
	var workingDir = t.TempDir()

	Expect(os.Mkdir(filepath.Join(workingDir, "custom"), os.ModePerm)).To(Succeed())

	// given a `pnpm-lock.yaml` exists and a `package.json` exists with a start script
	Expect(os.WriteFile(filepath.Join(workingDir, "custom", "pnpm-lock.yaml"), []byte{}, 0600)).To(Succeed())
	Expect(os.WriteFile(filepath.Join(workingDir, "custom", "package.json"), []byte(`{
				"scripts": {
					"start": "node server.js"
				}
			}`), 0600)).To(Succeed())

	// when detect is called
	result, err := detect(packit.DetectContext{
		WorkingDir: filepath.Join(workingDir, "custom"),
	})

	// then detect returns a successful build plan
	Expect(err).NotTo(HaveOccurred())
	Expect(result.Plan).To(Equal(packit.BuildPlan{
		Provides: []packit.BuildPlanProvision{
			{
				Name: NodeModules,
			},
		},
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
		},
	}))
}

func Test_DetectFailsWithoutStartScript(t *testing.T) {
	var Expect = NewWithT(t).Expect
	var workingDir = t.TempDir()

	Expect(os.Mkdir(filepath.Join(workingDir, "custom"), os.ModePerm)).To(Succeed())

	// given a `pnpm-lock.yaml` exists and a package.json exists WITHOUT a start script
	Expect(os.WriteFile(filepath.Join(workingDir, "custom", "pnpm-lock.yaml"), []byte{}, 0600)).To(Succeed())
	Expect(os.WriteFile(filepath.Join(workingDir, "custom", "package.json"), []byte(`{
				"scripts": {
					"lint":  "npm run lint"
				}
			}`), 0600)).To(Succeed())

	// when detect is called
	_, err := detect(packit.DetectContext{
		WorkingDir: filepath.Join(workingDir, "custom"),
	})

	// then detect fails with appropriate message
	Expect(err).To(MatchError(ContainSubstring("'package.json' has been found but does not have a 'start' command")))

	// cleanup
	Expect(os.RemoveAll(workingDir)).To(Succeed())
}

func Test_DetectFailsNoPackageJson(t *testing.T) {
	var Expect = NewWithT(t).Expect
	var workingDir = t.TempDir()

	Expect(os.Mkdir(filepath.Join(workingDir, "custom"), os.ModePerm)).To(Succeed())

	// given a `pnpm-lock/yaml` but no package json
	Expect(os.WriteFile(filepath.Join(workingDir, "custom", "pnpm-lock.yaml"), []byte{}, 0600)).To(Succeed())

	// when detect is called
	_, err := detect(packit.DetectContext{
		WorkingDir: filepath.Join(workingDir, "custom"),
	})

	// then detect fails with missing package.json message
	Expect(err).To(MatchError(packit.Fail.WithMessage("no 'package.json' found in project path %s", filepath.Join(workingDir, "custom"))))
}

func Test_DetectFailsNoPnpmLock(t *testing.T) {
	var Expect = NewWithT(t).Expect
	var workingDir = t.TempDir()

	Expect(os.Mkdir(filepath.Join(workingDir, "custom"), os.ModePerm)).To(Succeed())

	// given no pnpm-lock.yaml
	Expect(os.WriteFile(filepath.Join(workingDir, "custom", "package.json"), []byte(`{
				"scripts": {
					"start": "node server.js"
				}
			}`), 0600)).To(Succeed())

	// when detect is called
	_, err := detect(packit.DetectContext{
		WorkingDir: filepath.Join(workingDir, "custom"),
	})

	// then detect fails with missing pnpm-lock.yaml message
	Expect(err).To(MatchError(packit.Fail.WithMessage("no 'pnpm-lock.yaml' file found in the project path %s", filepath.Join(workingDir, "custom"))))
}

func Test_DetectFailsWithoutPackageJsonAccess(t *testing.T) {
	var Expect = NewWithT(t).Expect

	// create fake working directory
	var workingDir = t.TempDir()
	Expect(os.Mkdir(filepath.Join(workingDir, "custom"), os.ModePerm)).To(Succeed())

	// given `pnpm-lock.yaml` is created and `package.json` is created but doesn't have correct permissions
	Expect(os.WriteFile(filepath.Join(workingDir, "custom", "pnpm-lock.yaml"), []byte{}, 0600)).To(Succeed())

	var packageJsonPath = filepath.Join(workingDir, "custom", "package.json")
	Expect(os.WriteFile(packageJsonPath, []byte{}, 0600)).To(Succeed())
	Expect(os.Chmod(packageJsonPath, 0000)).To(Succeed())

	// when detect is called
	_, err := detect(packit.DetectContext{
		WorkingDir: filepath.Join(workingDir, "custom"),
	})

	// then detect fails with permission denied message
	Expect(err).To(MatchError(ContainSubstring("permission denied")))

	// cleanup
	Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
}
