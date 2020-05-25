package main

import (
	"github.com/dmikusa/rust-dist-cnb/rust"
	"github.com/paketo-buildpacks/packit"
)

func main() {
	packit.Detect(rust.Detect())
}
