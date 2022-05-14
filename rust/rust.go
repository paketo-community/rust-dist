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

func NewRust(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) Rust {
	expected := map[string]interface{}{"dependency": dependency}

	contributor := libpak.NewDependencyLayerContributor(
		dependency,
		cache,
		libcnb.LayerTypes{
			Build: true,
			Cache: true,
		})
	contributor.ExpectedMetadata = expected

	return Rust{
		LayerContributor: contributor,
	}
}

func (j Rust) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Expanding to %s", layer.Path)
		if err := crush.Extract(artifact, layer.Path, 1); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand Rust\n%w", err)
		}

		return layer, nil
	})
}

func (j Rust) Name() string {
	return j.LayerContributor.LayerName()
}
