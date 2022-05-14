package rust

import (
	"fmt"
	"os"
	"strings"

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

		// make layer for cargo, which is installed by rust
		cargo := Cargo{}
		cargo.Logger = b.Logger
		result.Layers = append(result.Layers, cargo)

		// install rust
		v, _ := cr.Resolve("BP_RUST_VERSION")

		rustDependency, err := dr.Resolve(PlanEntryRust, v)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		rust := NewRust(rustDependency, dc)
		rust.Logger = b.Logger
		result.Layers = append(result.Layers, rust)
	}

	return result, nil
}

func AppendToPath(values ...string) error {
	var path []string
	if curPath, ok := os.LookupEnv("PATH"); ok {
		path = append(path, curPath)
	}
	path = append(path, values...)
	return os.Setenv("PATH", strings.Join(path, string(os.PathListSeparator)))
}
