package rust

import (
	"github.com/buildpacks/libcnb"
)

// PlanEntryRust is the name of the plan
const PlanEntryRust = "rust"

type Detect struct {
}

// Detect will always just offer (i.e. provides) Rust, but does not require it ever
func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: PlanEntryRust},
				},
			},
		},
	}, nil
}
