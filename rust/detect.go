package rust

import (
	"os"

	"github.com/paketo-buildpacks/packit"
)

// PlanDependencyRust is the name of the plan
const PlanDependencyRust = "rust"

// BuildPlanMetadata defines the format of what is stored in the build plan's metadata section
type BuildPlanMetadata struct {
	Version       string `toml:"version"`
	VersionSource string `toml:"version-source"`
}

// Detect will always just offer (i.e. provides) Rust, but does not require it ever
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		var requirements []packit.BuildPlanRequirement

		version := os.Getenv("BP_RUST_VERSION")
		if version != "" {
			requirements = append(requirements, packit.BuildPlanRequirement{
				Name: PlanDependencyRust,
				Metadata: BuildPlanMetadata{
					Version:       version,
					VersionSource: "BP_RUST_VERSION",
				},
			})
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: PlanDependencyRust},
				},
				Requires: requirements,
			},
		}, nil
	}
}
