package rust_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dmikusa/rust-dist-cnb/rust"
	"github.com/dmikusa/rust-dist-cnb/rust/mocks"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		layersDir  string
		cnbPath    string
		timestamp  string

		dependencyService *mocks.DependencyService
		mockRunner        *mocks.Runner
		clock             rust.Clock

		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		layersDir, err = ioutil.TempDir("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbPath, err = ioutil.TempDir("", "cnb-path")
		Expect(err).NotTo(HaveOccurred())

		dependencyService = &mocks.DependencyService{}

		now := time.Now()
		clock = rust.NewClock(func() time.Time { return now })
		timestamp = now.Format(time.RFC3339Nano)

		logEmitter := rust.NewLogEmitter(ioutil.Discard)
		mockRunner = &mocks.Runner{}

		build = rust.Build(dependencyService, mockRunner, clock, logEmitter)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbPath)).To(Succeed())
	})

	it("installs rust", func() {
		dep := postal.Dependency{
			ID:           "rust",
			SHA256:       "some-sha",
			Source:       "some-source",
			SourceSHA256: "some-source-sha",
			Stacks:       []string{"some-stack"},
			URI:          "some-uri",
			Version:      "1.43.1",
		}
		dependencyService.On(
			"Resolve",
			mock.MatchedBy(func(s string) bool {
				return strings.HasSuffix(s, "buildpack.toml")
			}),
			"rust", "*", "some-stack",
		).Return(dep, nil)
		dependencyService.On("Install", dep, cnbPath, filepath.Join(layersDir, "downloads")).Return(nil)
		mockRunner.On(
			"Install",
			mock.MatchedBy(func(s string) bool {
				return strings.HasSuffix(s, "downloads")
			}),
			mock.MatchedBy(func(s string) bool {
				return strings.HasSuffix(s, "rust")
			}),
			"1.43.1",
		).Return(nil)

		result, err := build(packit.BuildContext{
			WorkingDir: workingDir,
			Layers:     packit.Layers{Path: layersDir},
			CNBPath:    cnbPath,
			Stack:      "some-stack",
			Plan: packit.BuildpackPlan{
				Entries: []packit.BuildpackPlanEntry{
					{
						Name: "rust",
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(packit.BuildResult{
			Layers: []packit.Layer{
				{
					Name:      "rust",
					Path:      filepath.Join(layersDir, "rust"),
					Build:     true,
					Cache:     true,
					SharedEnv: packit.Environment{},
					BuildEnv:  packit.Environment{},
					LaunchEnv: packit.Environment{},
					Metadata: map[string]interface{}{
						"built_at":  timestamp,
						"cache_sha": "some-sha",
					},
				},
			},
		}))
	})

	context("when rust was previously installed", func() {
		it.Before(func() {
			Expect(ioutil.WriteFile(filepath.Join(layersDir, "rust.toml"), []byte("launch = false\nbuild = true\ncache = true\n\n[metadata]\ncache_sha = \"some-sha\"\nbuilt_at = \"some_time\""), 0644)).To(Succeed())
		})

		it("skips the rust install", func() {
			dep := postal.Dependency{
				ID:           "rust",
				SHA256:       "some-sha",
				Source:       "some-source",
				SourceSHA256: "some-source-sha",
				Stacks:       []string{"some-stack"},
				URI:          "some-uri",
				Version:      "1.43.1",
			}
			dependencyService.On(
				"Resolve",
				mock.MatchedBy(func(s string) bool {
					return strings.HasSuffix(s, "buildpack.toml")
				}),
				"rust", "*", "some-stack",
			).Return(dep, nil)

			result, err := build(packit.BuildContext{
				WorkingDir: workingDir,
				Layers:     packit.Layers{Path: layersDir},
				CNBPath:    cnbPath,
				Stack:      "some-stack",
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name: "rust",
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(packit.BuildResult{
				Layers: []packit.Layer{
					{
						Name:      "rust",
						Path:      filepath.Join(layersDir, "rust"),
						Build:     true,
						Cache:     true,
						SharedEnv: packit.Environment{},
						BuildEnv:  packit.Environment{},
						LaunchEnv: packit.Environment{},
						Metadata: map[string]interface{}{
							"built_at":  "some_time",
							"cache_sha": "some-sha",
						},
					},
				},
			}))
		})
	})

	context("when the entry contains a version constraint", func() {
		it("builds rust with that version", func() {
			dep := postal.Dependency{
				ID:           "rust",
				SHA256:       "some-sha",
				Source:       "some-source",
				SourceSHA256: "some-source-sha",
				Stacks:       []string{"some-stack"},
				URI:          "some-uri",
				Version:      "1.43.1",
			}
			dependencyService.On(
				"Resolve",
				mock.MatchedBy(func(s string) bool {
					return strings.HasSuffix(s, "buildpack.toml")
				}),
				"rust", "1.43.1", "some-stack",
			).Return(dep, nil)
			dependencyService.On("Install", dep, cnbPath, filepath.Join(layersDir, "downloads")).Return(nil)
			mockRunner.On(
				"Install",
				mock.MatchedBy(func(s string) bool {
					return strings.HasSuffix(s, "downloads")
				}),
				mock.MatchedBy(func(s string) bool {
					return strings.HasSuffix(s, "rust")
				}),
				"1.43.1",
			).Return(nil)

			result, err := build(packit.BuildContext{
				WorkingDir: workingDir,
				Layers:     packit.Layers{Path: layersDir},
				CNBPath:    cnbPath,
				Stack:      "some-stack",
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name:    "rust",
							Version: "1.43.1",
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(packit.BuildResult{
				Layers: []packit.Layer{
					{
						Name:      "rust",
						Path:      filepath.Join(layersDir, "rust"),
						Build:     true,
						Cache:     true,
						SharedEnv: packit.Environment{},
						BuildEnv:  packit.Environment{},
						LaunchEnv: packit.Environment{},
						Metadata: map[string]interface{}{
							"built_at":  timestamp,
							"cache_sha": "some-sha",
						},
					},
				},
			}))
		})
	})

	context("failure cases", func() {
		context("when the rust layer cannot be retrieved", func() {
			it.Before(func() {
				Expect(ioutil.WriteFile(filepath.Join(layersDir, "rust.toml"), nil, 0000)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					Layers:     packit.Layers{Path: layersDir},
					CNBPath:    cnbPath,
					Stack:      "some-stack",
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{Name: "rust"},
						},
					},
				})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when the dependency cannot be resolved", func() {
			it.Before(func() {
				dependencyService.On(
					"Resolve",
					mock.MatchedBy(func(s string) bool {
						return strings.HasSuffix(s, "buildpack.toml")
					}),
					"rust", "*", "some-stack",
				).Return(postal.Dependency{}, errors.New("failed to resolve dependency"))
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					Layers:     packit.Layers{Path: layersDir},
					CNBPath:    cnbPath,
					Stack:      "some-stack",
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{Name: "rust"},
						},
					},
				})
				Expect(err).To(MatchError("failed to resolve dependency"))
			})
		})

		context("when the dependency cannot be installed", func() {
			it.Before(func() {
				dependencyService.On(
					"Resolve",
					mock.MatchedBy(func(s string) bool {
						return strings.HasSuffix(s, "buildpack.toml")
					}),
					"rust", "*", "some-stack",
				).Return(postal.Dependency{}, errors.New("failed to install dependency"))
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					Layers:     packit.Layers{Path: layersDir},
					CNBPath:    cnbPath,
					Stack:      "some-stack",
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{Name: "rust"},
						},
					},
				})
				Expect(err).To(MatchError("failed to install dependency"))
			})
		})
	})
}
