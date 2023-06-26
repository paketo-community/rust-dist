package rust

import (
	"fmt"
	"os"
	"path"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
	"github.com/paketo-buildpacks/libpak/effect"
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
		tempDir, err := os.MkdirTemp("", "rust")
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to create temporary directory to expand Rust\n%w", err)
		}
		if err := crush.Extract(artifact, tempDir, 1); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand Rust\n%w", err)
		}
		executor := effect.NewExecutor()
		j.Logger.Bodyf("Expanding to %s", layer.Path)
		if err = executor.Execute(effect.Execution{
			Command: "./install.sh",
			Args:    []string{fmt.Sprintf("--prefix=%s", layer.Path), "--disable-ldconfig"},
			Dir:     tempDir,
			Stdout:  bard.NewWriter(j.Logger.Logger.InfoWriter(), bard.WithIndent(3)),
			Stderr:  bard.NewWriter(j.Logger.Logger.InfoWriter(), bard.WithIndent(3)),
		}); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to install rust: %w", err)
		}
		if err = os.RemoveAll(path.Join(layer.Path, "share")); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to trim rust installation: %w", err)
		}
		return layer, nil
	})
}

func (j Rust) Name() string {
	return j.LayerContributor.LayerName()
}
