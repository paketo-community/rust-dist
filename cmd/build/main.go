package main

import (
	"os"
	"time"

	"github.com/dmikusa/rust-dist-cnb/rust"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/postal"
)

func main() {
	transport := cargo.NewTransport()
	dependencyService := postal.NewService(transport)
	clock := rust.NewClock(time.Now)
	logEmitter := rust.NewLogEmitter(os.Stdout)
	runner := rust.NewInstallRunner(pexec.NewExecutable("install.sh"))

	packit.Build(rust.Build(dependencyService, &runner, clock, logEmitter))
}
