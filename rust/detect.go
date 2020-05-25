package rust

import (
	"github.com/paketo-buildpacks/packit"
)

// PlanDependencyRust is the name of the plan
const PlanDependencyRust = "rust"

// Detect if the Rust binaries should be delivered
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: PlanDependencyRust},
				},
				Requires: []packit.BuildPlanRequirement{
					{Name: PlanDependencyRust},
				},
			},
		}, nil
	}
}
