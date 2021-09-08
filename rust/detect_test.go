package rust_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-community/rust-dist/rust"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect rust.Detect
	)

	it("includes build plan options", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "rust"},
					},
				},
			},
		}))
	})
}
