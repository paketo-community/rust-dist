package rust

import (
	"github.com/paketo-buildpacks/packit"
)

// PlanDependencyRust is the name of the plan
const PlanDependencyRust = "rust"

// Detect will always just offer (i.e. provides) Rust, but does not require it ever
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: PlanDependencyRust},
				},
			},
		}, nil
	}
}
