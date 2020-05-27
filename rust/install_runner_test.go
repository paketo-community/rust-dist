package rust_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/dmikusa/rust-dist-cnb/rust"
	"github.com/dmikusa/rust-dist-cnb/rust/mocks"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testInstallRunner(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect             = NewWithT(t).Expect
		downloadDir string = "/downloads"
		destDir     string = "/dest"
		runner      rust.InstallerRunner
		version     string = "1.2.3"
	)

	it.Before(func() {
		executable := &mocks.Executable{}
		execution := pexec.Execution{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Dir:    filepath.Join(downloadDir, fmt.Sprintf("rust-%s-x86_64-unknown-linux-gnu", version)),
			Args: []string{
				fmt.Sprintf("--prefix=%s", destDir),
				"--without=rust-docs,rls-preview,clippy-preview,miri-preview,rustfmt-preview,llvm-tools-preview,rust-analysis-x86_64-unknown-linux-gnu",
				"--disable-ldconfig",
			},
		}
		executable.On("Execute", execution).Return(nil)
		runner = rust.NewInstallRunner(executable, os.Stdout, os.Stderr)
	})

	it("runs the installer", func() {
		err := runner.Install(downloadDir, destDir, version)
		Expect(err).ToNot(HaveOccurred())
	})

	context("failure cases", func() {
		context("when the rust installer fails", func() {
			it.Before(func() {
				executable := &mocks.Executable{}
				execution := pexec.Execution{
					Stdout: os.Stdout,
					Stderr: os.Stderr,
					Dir:    filepath.Join(downloadDir, fmt.Sprintf("rust-%s-x86_64-unknown-linux-gnu", version)),
					Args: []string{
						fmt.Sprintf("--prefix=%s", destDir),
						"--without=rust-docs,rls-preview,clippy-preview,miri-preview,rustfmt-preview,llvm-tools-preview,rust-analysis-x86_64-unknown-linux-gnu",
						"--disable-ldconfig",
					},
				}
				executable.On("Execute", execution).Return(fmt.Errorf("expected"))
				runner = rust.NewInstallRunner(executable, os.Stdout, os.Stderr)
			})

			it("the error bubbles up", func() {
				err := runner.Install(downloadDir, destDir, version)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(Equal("expected")))
			})
		})
	})
}
