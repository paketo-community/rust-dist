package main

import (
	"os"

	"github.com/dmikusa/rust-dist-cnb/rust"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/draft"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
)

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout)
	entryResolver := draft.NewPlanner()
	dependencyService := postal.NewService(cargo.NewTransport())

	packit.Run(
		rust.Detect(),
		rust.Build(entryResolver, dependencyService, chronos.DefaultClock, logEmitter),
	)
}
