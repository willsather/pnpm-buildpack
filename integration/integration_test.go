package integration

import (
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func Test_SampleAppIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	Expect := NewWithT(t).Expect

	spec.Run(t, "Integration", func(t *testing.T, when spec.G, it spec.S) {
		var (
			buildpackPath string
			appPath       string
			output        []byte
			err           error
		)

		it.Before(func() {
			// Set up paths for the buildpack and sample app
			buildpackPath, err = filepath.Abs("../../bin/buildpack")
			Expect(err).NotTo(HaveOccurred())

			appPath, err = filepath.Abs("./sample-app")
			Expect(err).NotTo(HaveOccurred())
		})

		it("builds the app using the custom pnpm buildpack", func() {
			// Run the pack build command
			cmd := exec.Command("pack", "build", "sample-app", "--buildpack", buildpackPath, "--path", appPath)
			output, err = cmd.CombinedOutput()

			Expect(err).NotTo(HaveOccurred(), string(output))
			Expect(string(output)).To(ContainSubstring("Successfully built image"))
		})

		it("runs the app", func() {
			// Run the built image
			cmd := exec.Command("docker", "run", "-d", "-p", "3000:3000", "sample-app")
			containerID, err := cmd.Output()

			Expect(err).NotTo(HaveOccurred())

			// Verify the app is running
			defer exec.Command("docker", "rm", "-f", string(containerID)).Run()

			resp, err := exec.Command("curl", "http://localhost:3000").Output()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(resp)).To(ContainSubstring("Hello, world!"))
		})
	}, spec.Report(report.Terminal{}))
}
