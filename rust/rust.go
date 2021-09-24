package rust

import (
	"fmt"
	"os"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type Rust struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewRust(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (Rust, libcnb.BOMEntry) {
	expected := map[string]interface{}{"dependency": dependency}

	contributor, be := libpak.NewDependencyLayer(
		dependency,
		cache,
		libcnb.LayerTypes{
			Build: true,
			Cache: true,
		})
	contributor.ExpectedMetadata = expected

	// This is a workaround. At the moment, there is no feasible way with pack to see the build BOM, you can only see
    //   the launch BOM. We are including this dependency in the launch BOM for now to workaround this limitation.
    // When https://github.com/buildpacks/pack/issues/1221 is resolved and one can easily access report.toml
    //   we should remove this workaround.
	be.Launch = true

	return Rust{
		LayerContributor: contributor,
	}, be
}

func (j Rust) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Expanding to %s", layer.Path)
		if err := crush.ExtractTarGz(artifact, layer.Path, 1); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand Rust\n%w", err)
		}

		return layer, nil
	})
}

func (j Rust) Name() string {
	return j.LayerContributor.LayerName()
}
