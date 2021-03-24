package rust_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/dmikusa/rust-dist-cnb/rust"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		detect     packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		detect = rust.Detect()
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("returns a plan that provides rust", func() {
		result, err := detect(packit.DetectContext{
			WorkingDir: workingDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: rust.PlanDependencyRust},
				},
			},
		}))
	})

	context("when BP_RUST_VERSION is set", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_RUST_VERSION", "some-version")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_RUST_VERSION")).To(Succeed())
		})

		it("returns a plan that provides rust and requires a specific version", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(packit.DetectResult{
				Plan: packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: rust.PlanDependencyRust},
					},
					Requires: []packit.BuildPlanRequirement{
						{
							Name: rust.PlanDependencyRust,
							Metadata: rust.BuildPlanMetadata{
								VersionSource: "BP_RUST_VERSION",
								Version:       "some-version",
							},
						},
					},
				},
			}))
		})
	}, spec.Sequential())
}
