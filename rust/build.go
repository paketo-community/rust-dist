package rust

import (
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
)

//go:generate mockery -name DependencyService -case=underscore

// DependencyService interface for resolving and installing dependencies
type DependencyService interface {
	Resolve(path, name, version, stack string) (postal.Dependency, error)
	Install(dependency postal.Dependency, cnbPath, layerPath string) error
}

//go:generate mockery -name Runner -case=underscore

// Runner is something capable of running Cargo
type Runner interface {
	Install(downloadDir string, destDir string, version string) error
}

// Build does the actual install of Rust
func Build(dependencies DependencyService, runner Runner, clock chronos.Clock, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Resolving Rust version")
		logger.Candidates(context.Plan.Entries)

		entry := context.Plan.Entries[0]

		// downloadLayer, err := context.Layers.Get("downloads")
		// if err != nil {
		// 	return packit.BuildResult{}, err
		// }

		rustLayer, err := context.Layers.Get("rust")
		if err != nil {
			return packit.BuildResult{}, err
		}

		version := "*"
		if v, ok := entry.Metadata["version"].(string); ok {
			version = v
		}

		dependency, err := dependencies.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), "rust", version, context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.SelectedDependency(entry, dependency, clock.Now())

		if sha, ok := rustLayer.Metadata["cache_sha"].(string); !ok || sha != dependency.SHA256 {
			logger.Break()
			logger.Process("Installing Rust %s", dependency.Version)

			rustLayer, err = rustLayer.Reset()
			if err != nil {
				return packit.BuildResult{}, err
			}

			rustLayer.Build = true
			rustLayer.Cache = true

			logger.Subprocess("Downloading and extracting Rust")
			then := clock.Now()
			err = dependencies.Install(dependency, context.CNBPath, rustLayer.Path)
			if err != nil {
				return packit.BuildResult{}, err
			}
			logger.Action("Completed in %s", time.Since(then).Round(time.Millisecond))
			logger.Break()

			// logger.Subprocess("Installing Rust")
			// then = clock.Now()
			// err = runner.Install(downloadLayer.Path, rustLayer.Path, dependency.Version)
			// if err != nil {
			// 	return packit.BuildResult{}, err
			// }
			// logger.Action("Completed in %s", time.Since(then).Round(time.Millisecond))
			// logger.Break()

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
