package rust_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-community/rust-dist/rust"
	"github.com/sclevine/spec"
)

func testRust(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		ctx    libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path, err = ioutil.TempDir("", "rust-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes Rust", func() {
		dep := libpak.BuildpackDependency{
			Version: "1.54.0",
			URI:     "https://localhost/stub-rust-1.54.0.tar.gz",
			SHA256:  "2b61f788c000462840d4bf0f83f5ad3d552a8d27f364dd5f7e7e6b29031bee43",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		j := rust.NewRust(dep, dc)
		j.Logger = bard.NewLogger(ioutil.Discard)

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = j.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.LayerTypes.Launch).To(BeFalse())
		Expect(layer.LayerTypes.Build).To(BeTrue())
		Expect(layer.LayerTypes.Cache).To(BeTrue())

		Expect(filepath.Join(layer.Path, "fixture-marker")).To(BeARegularFile())
		Expect(filepath.Join(layer.Path, "bin")).To(BeADirectory())
		Expect(layer.SBOMPath(libcnb.SyftJSON)).To(BeARegularFile())
	})
}
