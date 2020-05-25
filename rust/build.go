package rust

import (
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/postal"
)

//go:generate mockery -name DependencyService -case=underscore

// DependencyService interface for resolving and installing dependencies
type DependencyService interface {
	Resolve(path, name, version, stack string) (postal.Dependency, error)
	Install(dependency postal.Dependency, cnbPath, layerPath string) error
}

// Build does the actual install of Rust
func Build(dependencies DependencyService, clock Clock, logger LogEmitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title(context.BuildpackInfo)

		logger.Process("Resolving Rust version")
		logger.Candidates(context.Plan.Entries)

		entry := context.Plan.Entries[0]

		rustLayer, err := context.Layers.Get("rust", packit.BuildLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		version := "*"
		if entry.Version != "" {
			version = entry.Version
		}

		dependency, err := dependencies.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), "rust", version, context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.SelectedDependency(entry, dependency.Version)

		if sha, ok := rustLayer.Metadata["cache_sha"].(string); !ok || sha != dependency.SHA256 {
			logger.Break()
			logger.Process("Executing build process")

			err = rustLayer.Reset()
			if err != nil {
				return packit.BuildResult{}, err
			}

			logger.Subprocess("Installing Rust %s", dependency.Version)
			then := clock.Now()
			err = dependencies.Install(dependency, context.CNBPath, rustLayer.Path)
			if err != nil {
				return packit.BuildResult{}, err
			}
			logger.Action("Completed in %s", time.Since(then).Round(time.Millisecond))
			logger.Break()

			rustLayer.Metadata = map[string]interface{}{
				"built_at":  clock.Now().Format(time.RFC3339Nano),
				"cache_sha": dependency.SHA256,
			}
		}

		return packit.BuildResult{
			Layers: []packit.Layer{
				rustLayer,
			},
		}, nil
	}
}
