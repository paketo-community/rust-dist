package rust

import (
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
)

//go:generate mockery -name EntryResolver -case=underscore

// EntryResolver for resolving buildpack plan entries
type EntryResolver interface {
	Resolve(name string, entries []packit.BuildpackPlanEntry, priorites []interface{}) (packit.BuildpackPlanEntry, []packit.BuildpackPlanEntry)
}

//go:generate mockery -name DependencyService -case=underscore

// DependencyService interface for resolving and installing dependencies
type DependencyService interface {
	Resolve(path, name, version, stack string) (postal.Dependency, error)
	Install(dependency postal.Dependency, cnbPath, layerPath string) error
}

// Priorities defines the order in which we select versions
var Priorities = []interface{}{
	"BP_RUST_VERSION",
	"CARGO",
}

// Build does the actual install of Rust
func Build(entryResolver EntryResolver, dependencies DependencyService, clock chronos.Clock, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Resolving Rust version")
		entry, entries := entryResolver.Resolve(PlanDependencyRust, context.Plan.Entries, Priorities)
		logger.Candidates(entries)

		version, ok := entry.Metadata["version"].(string)
		if !ok {
			version = "default"
		}

		dependency, err := dependencies.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), "rust", version, context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.SelectedDependency(entry, dependency, clock.Now())

		rustLayer, err := context.Layers.Get("rust")
		if err != nil {
			return packit.BuildResult{}, err
		}

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
