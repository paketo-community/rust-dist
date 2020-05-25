package rust_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

	it("returns a DetectResult that provides rust", func() {
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

	context("when there is a *.rs file in the workspace", func() {
		it.Before(func() {
			_, err := os.Create(filepath.Join(workingDir, "src", "lib.rs"))
			Expect(err).NotTo(HaveOccurred())
		})

		it("returns a DetectResult that provides and required rust", func() {
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
						},
					},
				},
			}))
		})

		// context("when the buildpack.yml specifies a version to install", func() {
		// 	it.Before(func() {
		// 		err := ioutil.WriteFile(filepath.Join(workingDir, "buildpack.yml"), []byte(`{"httpd": {"version": "1.2.3"}}`), 0644)
		// 		Expect(err).NotTo(HaveOccurred())
		// 	})

		// 	it("returns a DetectResult that provides and required httpd with that version", func() {
		// 		result, err := detect(packit.DetectContext{
		// 			WorkingDir: workingDir,
		// 		})
		// 		Expect(err).NotTo(HaveOccurred())
		// 		Expect(result).To(Equal(packit.DetectResult{
		// 			Plan: packit.BuildPlan{
		// 				Provides: []packit.BuildPlanProvision{
		// 					{Name: rust.PlanDependencyRust},
		// 				},
		// 				Requires: []packit.BuildPlanRequirement{
		// 					{
		// 						Name:    rust.PlanDependencyRust,
		// 						Version: "1.2.3",
		// 						Metadata: rust.BuildPlanMetadata{
		// 							Launch:        true,
		// 							VersionSource: "buildpack.yml",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		}))
		// 	})
		// })
	})

	// context("failure cases", func() {
	// context("when the buildpack.yml is malformed", func() {
	// 	it.Before(func() {
	// 		err := ioutil.WriteFile(filepath.Join(workingDir, "buildpack.yml"), []byte("%%%"), 0644)
	// 		Expect(err).NotTo(HaveOccurred())
	// 	})

	// 	it("returns an error", func() {
	// 		_, err := detect(packit.DetectContext{WorkingDir: workingDir})
	// 		Expect(err).To(MatchError(ContainSubstring("failed to parse buildpack.yml")))
	// 	})
	// })

	// context("when there is a buildpack.yml without a *.rs file", func() {
	// 	it.Before(func() {
	// 		err := ioutil.WriteFile(filepath.Join(workingDir, "buildpack.yml"), []byte(`{"rust": {"version": "1.2.3"}}`), 0644)
	// 		Expect(err).NotTo(HaveOccurred())
	// 	})

	// 	it("returns an error", func() {
	// 		_, err := detect(packit.DetectContext{WorkingDir: workingDir})
	// 		Expect(err).To(MatchError("failed to detect: buildpack.yml specifies a version, but no rust files present"))
	// 	})
	// })
	// })
}
