package rust

import (
	"io"
	"strconv"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

// LogEmitter to be used by the CNB
type LogEmitter struct {
	scribe.Logger
}

// NewLogEmitter makes a new pretty CNB logger
func NewLogEmitter(output io.Writer) LogEmitter {
	return LogEmitter{
		Logger: scribe.NewLogger(output),
	}
}

// Title is for the initial buildpack title
func (e LogEmitter) Title(info packit.BuildpackInfo) {
	e.Logger.Title("%s %s", info.Name, info.Version)
}

// Candidates logs a pretty list of dependency version candidates
func (e LogEmitter) Candidates(entries []packit.BuildpackPlanEntry) {
	e.Logger.Subprocess("Candidate version sources (in priority order):")

	var (
		sources [][2]string
		maxLen  int
	)

	for _, entry := range entries {
		versionSource, ok := entry.Metadata["version-source"].(string)
		if !ok {
			versionSource = "<unknown>"
		}

		if len(versionSource) > maxLen {
			maxLen = len(versionSource)
		}

		sources = append(sources, [2]string{versionSource, entry.Version})
	}

	for _, source := range sources {
		e.Logger.Action(("%-" + strconv.Itoa(maxLen) + "s -> %q"), source[0], source[1])
	}

	e.Logger.Break()
}

// SelectedDependency logs a pretty version of what dependency was picked
func (e LogEmitter) SelectedDependency(entry packit.BuildpackPlanEntry, version string) {
	source, ok := entry.Metadata["version-source"].(string)
	if !ok {
		source = "<unknown>"
	}

	e.Logger.Subprocess("Selected Rust version (using %s): %s", source, version)
}

// Environment logs environment values set by the CNB
func (e LogEmitter) Environment(environment packit.Environment) {
	e.Logger.Subprocess("%s", scribe.NewFormattedMapFromEnvironment(environment))
}
