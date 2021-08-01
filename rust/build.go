package rust

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	if _, ok, err := pr.Resolve(PlanEntryRust); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve Rust plan entry\n%w", err)
	} else if ok {
		dc, err := libpak.NewDependencyCache(context)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
		}
		dc.Logger = b.Logger

		dr, err := libpak.NewDependencyResolver(context)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
		}

		cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
		}

		v, _ := cr.Resolve("BP_RUST_VERSION")

		rustupDependency, err := dr.Resolve(PlanEntryRust, v)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		rustup, be := NewRust(rustupDependency, dc)
		rustup.Logger = b.Logger

		result.Layers = append(result.Layers, rustup)
		result.BOM.Entries = append(result.BOM.Entries, be)
	}

	return result, nil
}
